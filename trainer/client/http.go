package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/trainer/config"
	"github.com/MingkaiLee/kasos/trainer/util"
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
