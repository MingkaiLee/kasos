package internal

import (
	"github.com/MingkaiLee/kasos/infer-module/util"
)

var InferCronJob *CronJob

// 每次启动时需要先拉取当前已注册的服务并创建CronJob
func InitInternal() {
	svcs, err := ListServices()
	if err != nil {
		// 如果拉取失败, 上报错误日志, 服务仍会启动
		util.LogErrorf("init internal error: %v", err)
	}
	parallelInferer := NewParallelInferer(svcs)
	InferCronJob = NewCronJob(15, parallelInferer)
	InferCronJob.Start()
}
