package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/infer-module/config"
	"github.com/MingkaiLee/kasos/infer-module/util"
	jsoniter "github.com/json-iterator/go"
)

var httpClient *http.Client

func InitHTTPClient() {
	// 建立一个http连接池
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     60 * time.Second,
		},
	}
}

func CallListHpaServices(ctx context.Context, index int) (response *http.Response, err error) {
	url := fmt.Sprintf("%s/%s?index=%d", config.ServerUrl, "service-manager/list", index)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		util.LogErrorf("http.CallListHpaServices error: %s", err)
		return
	}
	response, err = httpClient.Do(request)

	return
}

func CallFindModel(ctx context.Context, modelName string) (response *http.Response, err error) {
	url := fmt.Sprintf("%s/%s?name=%s", config.ServerUrl, "model-manager/find", modelName)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		util.LogErrorf("http.NewRequest error: %s", err)
		return
	}
	response, err = httpClient.Do(request)

	return
}

type ReportQPSRequest struct {
	ServiceName string `json:"service_name"`
	QPS         int    `json:"qps"`
}

func CallReportQPS(ctx context.Context, req *ReportQPSRequest) (response *http.Response, err error) {
	if req == nil {
		err = fmt.Errorf("request is nil")
		util.LogErrorf("http.CallReportQPS error: %v", err)
		return
	}
	url := fmt.Sprintf("%s/%s", config.HpaExecutorUrl, "hpa-exec/report-qps")
	content, err := jsoniter.Marshal(*req)
	if err != nil {
		util.LogErrorf("http.CallReportThresh error: %v", err)
		return
	}
	body := bytes.NewReader(content)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		util.LogErrorf("http.CallReportThresh error: %s", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = httpClient.Do(request)
	if err != nil {
		util.LogErrorf("http.CallReportThresh error: %s", err)
		return
	}

	return
}

type ReportModelValidRequest struct {
	ModelName string `json:"model_name"`
	Ok        bool   `json:"ok"`
	ErrorInfo string `json:"error_info"`
}

func CallReportModelValid(ctx context.Context, req *ReportModelValidRequest) (response *http.Response, err error) {
	if req == nil {
		err = fmt.Errorf("request is nil")
		util.LogErrorf("http.CallReportModelValid error: %v", err)
		return
	}
	url := fmt.Sprintf("%s/%s", config.ServerUrl, "model-manager/report-valid")
	content, err := jsoniter.Marshal(*req)
	if err != nil {
		util.LogErrorf("http.CallReportModelValid error: %v", err)
		return
	}
	body := bytes.NewReader(content)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		util.LogErrorf("http.CallReportModelValid error: %v", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err = httpClient.Do(request)
	if err != nil {
		util.LogErrorf("http.CallReportModelValid error: %v", err)
		return
	}

	return
}
