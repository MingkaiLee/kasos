package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MingkaiLee/kasos/hpa-executor/config"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"
)

var validHTTPMethods []string = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodPut,
}

type NormalTester struct {
	client      *http.Client           // 进行压测的HTTP客户端
	name        string                 // 压测服务名
	result      chan *StressTestResult // 压测结果
	request     *http.Request          // 压测请求内容
	method      string                 // 压测的HTTP方法
	url         string                 // 压测目标的URL
	contentType string                 // 压测目标所需要的body填充类型
	content     string                 // 压测目标所需要的body
	initialQPS  int64                  // 压测的起步QPS
	rt          int64                  // 接口的基本延时, ms
}

type NormalTesterSettings struct {
	Name        string `json:"name"`
	Method      string `json:"method"`
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	InitialQPS  *int64 `json:"initial_qps"`
	Timeout     *int64 `json:"timeout"`
}

func NewNormalTester() *NormalTester {
	return &NormalTester{
		client: &http.Client{Timeout: time.Duration(config.DefaultClientTimeout) * time.Millisecond},
		result: make(chan *StressTestResult, 1),
	}
}

// 初始化配置
func (p *NormalTester) SetConfig(settings *NormalTesterSettings) (err error) {
	ok := slices.Contains(validHTTPMethods, settings.Method)
	if !ok {
		err = fmt.Errorf("unknown http method: %s, only support: GET, POST, DELETE, PUT", settings.Method)
		util.LogErrorf("NormalTester.SetConfig error: %v", err)
		return
	}

	p.name = settings.Name
	p.method = settings.Method
	p.url = settings.Url
	p.contentType = settings.ContentType
	p.content = settings.Content
	// 测试一次看访问方式是否可行
	req, err := p.genReq()
	if err != nil {
		util.LogErrorf("NormalTester.SetConfig error: %v", err)
		return
	}
	_, err = p.client.Do(req)
	if err != nil {
		util.LogErrorf("NormalTester.SetConfig error: %v", err)
		return
	}
	p.request = req

	if settings.InitialQPS != nil {
		p.initialQPS = *settings.InitialQPS
	}
	// 设置超时时间
	if settings.Timeout != nil {
		p.client.Timeout = time.Duration(*settings.Timeout) * time.Second
	}

	return
}

// 通过JSON初始化配置
func (p *NormalTester) SetConfigByJSON(json []byte) (err error) {
	settings := NormalTesterSettings{}
	err = jsoniter.Unmarshal(json, &settings)
	if err != nil {
		util.LogErrorf("NormalTester.SetConfigByJSON, error: %v", err)
		return
	}
	err = p.SetConfig(&settings)
	if err != nil {
		util.LogErrorf("NormalTester.SetConfigByJSON, error: %v", err)
	}
	return
}

func (p *NormalTester) genReq() (*http.Request, error) {
	var body io.Reader
	if p.content != "" {
		body = strings.NewReader(p.content)
	}
	req, err := http.NewRequest(p.method, p.url, body)
	if err != nil {
		return nil, err
	}
	if p.contentType != "" {
		req.Header.Set("Content-Type", p.contentType)
	}
	return req, nil
}

// 测出一个接口的基本延时
func (p *NormalTester) basicRtTest() {
	var totalTime int64 = 0
	var reqCount int64 = 0
	for i := 0; i < config.RtTestEpoch; i++ {
		startTime := time.Now().UnixMicro()
		_, err := p.client.Do(p.request)
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

// 开启goroutine压测
func (p *NormalTester) Run() {
	go func() {
		r := new(StressTestResult)
		r.ServiceName = p.name
		// 将结果送入channel中
		defer func() {
			util.LogInfof("NormalTester.Run, test finished, result: %+v", r)
			p.result <- r
		}()
		// 判断是否已完成初始化
		if p.request == nil {
			err := fmt.Errorf("request is nil, please call SetConfigByJSON first")
			util.LogErrorf("NormalTester.Run error: %v", err)
			r.Err = err
			return
		}
		// 基本延时测试
		p.basicRtTest()
		// 填充接口的基本延迟
		r.Rt = p.rt
		// 开始压测
		for qps := p.initialQPS; qps < config.MaxQPS; qps++ {
			// 依靠WaitGroup和goroutine并发请求
			var m sync.WaitGroup
			var successCount int64 = 0
			var totalRt int64 = 0
			var idx int64
			for idx = 0; idx < qps; idx++ {
				m.Add(1)
				go func() {
					defer m.Done()
					startTime := time.Now().UnixMilli()
					_, err := p.client.Do(p.request)
					if err != nil {
						return
					}
					endTime := time.Now().UnixMilli()
					diff := endTime - startTime
					// 使用atomic保证并发安全
					atomic.AddInt64(&successCount, 1)
					atomic.AddInt64(&totalRt, diff)
				}()
			}
			m.Wait()
			// 是否已达到阈值标志
			threshReached := false
			// 检查rt
			avgRt := totalRt / successCount
			if avgRt > 2*p.rt {
				util.LogInfof("NormalTester.Run, the rt reached the limit")
				threshReached = true
			}
			// 检查错误率
			errorRate := float64(qps-successCount) / float64(qps)
			if errorRate > config.ErrorTolerateRate {
				util.LogInfof("NormalTester.Run, the error rate reached the limit")
				threshReached = true
			}
			// 达到阈值, 退出循环
			if threshReached {
				// 填充测试完成的阈值QPS
				r.ThresholdQPS = int(qps)
				break
			}
		}
	}()
}

// 压测结果的channel
func (p *NormalTester) GetResult() <-chan *StressTestResult {
	return p.result
}
