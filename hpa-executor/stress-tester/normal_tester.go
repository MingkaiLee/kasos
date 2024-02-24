package stresstester

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"

	"github.com/MingkaiLee/kasos/hpa-executor/conf"
)

var validHTTPMethods []string = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodPut,
}

type NormalTester struct {
	client      *http.Client
	result      *StressTestResult
	requst      *http.Request
	method      string
	url         string
	contentType string
	content     string
	initialQPS  int64
	rt          int64
}

type NormalTesterSettings struct {
	Method      string `json:"method"`
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	InitialQPS  *int64 `json:"initial_qps"`
	Timeout     *int64 `json:"timeout"`
}

func NewNormalTester() *NormalTester {
	return &NormalTester{
		client: &http.Client{Timeout: time.Duration(conf.NormalTesterConf.DefaultClientTimeout) * time.Second},
	}
}

func (p *NormalTester) SetConfigByJSON(json string) error {
	settings := NormalTesterSettings{}
	err := jsoniter.UnmarshalFromString(json, &settings)
	if err != nil {
		return err
	}

	ok := slices.Contains(validHTTPMethods, settings.Method)
	if !ok {
		return fmt.Errorf("unknown http method: %s, only support: GET, POST, DELETE, PUT", settings.Method)
	}

	p.method = settings.Method
	p.url = settings.Url
	p.contentType = settings.ContentType
	p.content = settings.Content
	// 测试一次看访问方式是否可行
	req, err := p.genReq()
	if err != nil {
		return err
	}
	_, err = p.client.Do(req)
	if err != nil {
		return err
	}
	p.requst = req

	if settings.InitialQPS != nil {
		p.initialQPS = *settings.InitialQPS
	}
	if settings.Timeout != nil {
		p.client.Timeout = time.Duration(*settings.Timeout) * time.Second
	}

	return nil
}

func (p *NormalTester) genReq() (*http.Request, error) {
	req, err := http.NewRequest(p.method, p.url, strings.NewReader(p.content))
	if err != nil {
		return nil, err
	}
	if p.contentType != "" {
		req.Header.Set("Content-Type", p.contentType)
	}
	return req, nil
}

func (p *NormalTester) basicRtTest() {
	var totalTime int64 = 0
	var reqCount int64 = 0
	for i := 0; i < conf.NormalTesterConf.RtTestEpoch; i++ {
		startTime := time.Now().UnixMicro()
		_, err := p.client.Do(p.requst)
		if err != nil {
			continue
		}
		endTime := time.Now().UnixMicro()
		diff := endTime - startTime
		totalTime += diff
		reqCount += 1
	}
	p.rt = totalTime / reqCount
}

func (p *NormalTester) Run() {}
