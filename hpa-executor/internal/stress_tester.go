package internal

import (
	"context"

	"github.com/MingkaiLee/kasos/hpa-executor/client"
	"github.com/MingkaiLee/kasos/hpa-executor/config"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
)

type StressTestResult struct {
	ServiceName  string // 服务名
	ThresholdQPS int    // 阈值QPS
	Rt           int64  // 接口响应时间, ms
	Err          error  // 错误信息
}

// 抽象压测接口
type StreeTester interface {
	SetConfigByJSON(json string) error
	Run()
	GetResult() <-chan *StressTestResult
}

// 开启异步任务压测
func NormalTest(ctx context.Context, conf []byte) (err error) {
	tester := NewNormalTester()
	err = tester.SetConfigByJSON(conf)
	if err != nil {
		util.LogErrorf("service.NormalTest error: %v", err)
		return
	}
	// 执行压测任务
	tester.Run()
	// 异步回调
	go func() {
		result := <-tester.GetResult()
		util.LogInfof("service.NormalTest result: %+v", result)
		var req client.ReportThreshRequest
		if result.Err != nil {
			req.ServiceName = result.ServiceName
			req.OK = false
			req.ErrorInfo = result.Err.Error()
			req.ThreshQPS = uint(result.ThresholdQPS)
		} else {
			req.ServiceName = result.ServiceName
			req.OK = true
			req.ThreshQPS = uint(result.ThresholdQPS)
			// 更新全局map
			config.UpDateServiceThreshQPSCache(req.ServiceName, int(req.ThreshQPS))
		}
		// 发送回调消息
		resp, e := client.CallReportThresh(ctx, &req)
		if e != nil {
			util.LogErrorf("service.NormalTest report request: %+v,resp: %+v, err: %v", req, *resp, e)
		}
		// 打印日志
		util.LogInfof("service.NormalTest report request: %+v,resp: %+v, err: %v", req, *resp, e)
	}()
	return
}
