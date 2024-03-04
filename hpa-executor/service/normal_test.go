package service

import (
	"context"

	"github.com/MingkaiLee/kasos/hpa-executor/client"
	"github.com/MingkaiLee/kasos/hpa-executor/internal"
	"github.com/MingkaiLee/kasos/hpa-executor/util"
)

func NormalTest(ctx context.Context, conf []byte) (err error) {
	tester := internal.NewNormalTester()
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
		}
		// 发送回调消息
		resp, e := client.CallReportThresh(ctx, &req)
		// 打印日志
		util.LogInfof("service.NormalTest report request: %+v,resp: %+v, err: %v", req, *resp, e)
	}()
	return
}
