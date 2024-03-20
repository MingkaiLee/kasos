package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MingkaiLee/kasos/server/client"
	"github.com/MingkaiLee/kasos/server/model"
	"github.com/MingkaiLee/kasos/server/util"
	jsoniter "github.com/json-iterator/go"
)

type RegisterServiceRequest struct {
	Name      string            `json:"name"`
	Tags      map[string]string `json:"tags"`
	ModelName string            `json:"model_name"`
}

type RegisterServiceResponse struct {
	Accepted bool   `json:"accepted"`
	Message  string `json:"message"`
}

func RegisterService(ctx context.Context, content []byte) (response *RegisterServiceResponse, err error) {
	var req RegisterServiceRequest
	err = jsoniter.Unmarshal(content, &req)
	response = new(RegisterServiceResponse)
	if err != nil {
		util.LogErrorf("failed to unmarshal request, error: %v", err)
		response.Message = err.Error()
		return
	}
	util.LogInfof("register service request: %+v", req)

	// look up model information
	_, err = model.HpaModelGetID(req.ModelName)
	if err != nil {
		util.LogErrorf("failed to get model id, error: %v", err)
		response.Message = err.Error()
		return
	}

	// register service
	err = model.HpaServiceCreate(req.Name, req.Tags, req.ModelName)
	if err != nil {
		util.LogErrorf("failed to create service, error: %v", err)
		response.Message = err.Error()
		return
	}

	// 如果失败了删除记录
	defer func(e error) {
		if e != nil {
			model.HpaServiceDelete(req.Name)
		}
	}(err)

	// send stress test request to hpa-executor
	// 固定测试的接口
	url := fmt.Sprintf("http://%s.default.svc.cluster.local:8080/stress-test", req.Name)
	// 默认从1开始压测
	var initialQPS int64 = 1
	// 默认10秒超时
	var timeout int64 = 10
	testConf := client.NormalTesterSettings{
		Name:       req.Name,
		Method:     "GET",
		Url:        url,
		InitialQPS: &initialQPS,
		Timeout:    &timeout,
	}
	testResp, err := client.CallNormalTest(ctx, &testConf)
	if err != nil {
		util.LogErrorf("service.RegisterService error: %v", err)
		response.Message = err.Error()
		return
	}
	if testResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("call tester failed, status: %v, code: %d", testResp.Status, testResp.StatusCode)
		util.LogErrorf("service.RegisterService error: %v", err)
		response.Message = err.Error()
		return
	}

	// create a Prometheus ServiceMonitor
	err = client.CreateMonitorService(ctx, req.Name, req.Tags)
	if err != nil {
		util.LogErrorf("service.RegisterService error: %v", err)
		response.Message = err.Error()
		return
	}

	// 注册服务到infer-module
	addServiceReq := client.AddServiceRequest{
		ServiceName: req.Name,
		ModelName:   req.ModelName,
		Tags:        util.ConvertTags(req.Tags),
	}
	addResp, err := client.CallAddService(ctx, &addServiceReq)
	if err != nil {
		util.LogErrorf("service.RegisterService error: %v", err)
		response.Message = err.Error()
		return
	}
	// 解析返回结果
	if addResp.StatusCode != http.StatusOK {
		err = fmt.Errorf("call add service failed, status: %v, code: %d", addResp.Status, addResp.StatusCode)
		util.LogErrorf("service.RegisterService error: %v", err)
		response.Message = err.Error()
		return
	}

	response.Accepted = true
	return
}
