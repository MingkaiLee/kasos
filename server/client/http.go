package client

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func InitHTTPClient() {
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     60 * time.Second,
		},
	}
}
