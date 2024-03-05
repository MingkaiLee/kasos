package internal

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/MingkaiLee/kasos/hpa-executor/client"
	"github.com/MingkaiLee/kasos/hpa-executor/config"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
	jsoniter "github.com/json-iterator/go"
)

type ReportQPSRequest struct {
	ServiceName string `json:"service_name"`
	QPS         int    `json:"qps"`
}

type HpaService struct {
	Name      *string           `json:"name"`
	Tags      map[string]string `json:"tags"`
	ThreshQPS *uint             `json:"thresh_qps"`
	ModelName *string           `json:"model_name"`
}

// 上报QPS预测值后, 执行扩缩容
func HpaExec(ctx context.Context, content []byte) (err error) {
	var req ReportQPSRequest
	err = jsoniter.Unmarshal(content, &req)
	if err != nil {
		util.LogErrorf("service.HpaExec error: %v", err)
		return
	}
	// 获取当前的QPS阈值
	threshQPS, ok := config.GetServiceThreshQPS(req.ServiceName)
	if !ok {
		util.LogInfof("get local thresh qps failed, service: %s", req.ServiceName)
		// 如果本地的QPS阈值不存在, 向管控模块查询
		resp, e := client.CallFindHpaService(ctx, req.ServiceName)
		if e != nil {
			err = e
			util.LogErrorf("service.HpaExec error: %v", err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unexpected response, status: %s, code:: %d", resp.Status, resp.StatusCode)
			util.LogErrorf("service.HpaExec error: %v", err)
			return
		}
		defer resp.Body.Close()
		body, e := io.ReadAll(resp.Body)
		if e != nil {
			err = e
			util.LogErrorf("service.HpaExec error: %v", err)
			return
		}
		var hpaService HpaService
		err = jsoniter.Unmarshal(body, &hpaService)
		if err != nil {
			util.LogErrorf("service.HpaExec error: %v", err)
			return
		}
		// 更新本地map
		util.LogInfof("update local thresh qps, service: %s, thresh qps: %d", req.ServiceName, *hpaService.ThreshQPS)
		config.UpDateServiceThreshQPSCache(req.ServiceName, int(*hpaService.ThreshQPS))
		threshQPS = int(*hpaService.ThreshQPS)
	}
	// 计算目标replica数量
	targetReplica := int32(req.QPS / threshQPS)
	if (req.QPS % threshQPS) != 0 {
		targetReplica++
	}
	// 获取当前的deployment配置
	d, err := client.GetDeployment(ctx, req.ServiceName)
	if err != nil {
		util.LogErrorf("service.HpaExec error: %v", err)
		return
	}
	util.LogInfof("update service replicas, service: %s, current replicas: %d, target replicas: %d",
		req.ServiceName,
		*d.Spec.Replicas,
		targetReplica)
	// 比较并更新
	if *d.Spec.Replicas != targetReplica {
		d.Spec.Replicas = &targetReplica
		e := client.UpdateDeplotment(ctx, d)
		if e != nil {
			util.LogErrorf("service.HpaExec error: %v", e)
			return
		}
		util.LogInfof("update service replicas success, service: %s", req.ServiceName)
	}
	return
}
