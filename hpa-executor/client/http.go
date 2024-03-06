package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/hpa-executor/config"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
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
	url := fmt.Sprintf("%s/%s?index=%d", config.ServerUrl, "/service-manager/list", index)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		util.LogErrorf("http.CallListHpaServices error: %s", err)
		return
	}
	response, err = httpClient.Do(request)

	return
}

func CallFindHpaService(ctx context.Context, serviceName string) (response *http.Response, err error) {
	url := fmt.Sprintf("%s/%s?name=%s", config.ServerUrl, "/service-manager/find", serviceName)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		util.LogErrorf("http.CallFindHpaService error: %v", err)
		return
	}
	response, err = httpClient.Do(request)

	return
}

type ReportThreshRequest struct {
	ServiceName string `json:"service_name"`
	OK          bool   `json:"ok"`
	ErrorInfo   string `json:"error_info"`
	ThreshQPS   uint   `json:"thresh_qps"`
}

func CallReportThresh(ctx context.Context, req *ReportThreshRequest) (response *http.Response, err error) {
	if req == nil {
		err = fmt.Errorf("request is nil")
		util.LogErrorf("http.CallReportThresh error: %v", err)
		return
	}
	url := fmt.Sprintf("%s/%s", config.ServerUrl, "/service-manager/report-thresh")
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
