package client

import (
	"net/http"
	"time"
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
