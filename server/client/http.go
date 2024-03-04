package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/server/config"
	"github.com/MingkaiLee/kasos/server/util"
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

// 调用压测模块的请求体
type NormalTesterSettings struct {
	Name        string `json:"name"`
	Method      string `json:"method"`
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	InitialQPS  *int64 `json:"initial_qps"`
	Timeout     *int64 `json:"timeout"`
}

func CallNormalTest(ctx context.Context, req *NormalTesterSettings) (response *http.Response, err error) {
	if req == nil {
		err = fmt.Errorf("request is nil")
		util.LogErrorf("http.CallNormalTest error: %v", err)
		return
	}
	url := fmt.Sprintf("%s/%s", config.HpaExecutorUrl, "/stress-test/normal-test")
	content, err := jsoniter.Marshal(*req)
	if err != nil {
		util.LogErrorf("http.CallNormalTest error: %v", err)
		return
	}
	body := bytes.NewReader(content)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		util.LogErrorf("http.CallNormalTest error: %s", err)
		return
	}
	response, err = httpClient.Do(request)
	if err != nil {
		util.LogErrorf("http.CallNormalTest error: %s", err)
		return
	}

	return
}
