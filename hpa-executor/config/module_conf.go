package config

import (
	"os"

	"github.com/MingkaiLee/kasos/hpa-executor/util"
	jsoniter "github.com/json-iterator/go"
)

var (
	DefaultClientTimeout int64
	RtTestEpoch          int
	MaxQPS               int64
	ErrorTolerateRate    float64
)

const moduleConf = "/etc/config/module.json"

type NormalTesterConf struct {
	DefaultClientTimeout int64   `json:"default_client_timeout"` // 默认的client超时时间, ms
	RtTestEpoch          int     `json:"rt_test_epoch"`          // rt测试的epoch
	MaxQPS               int64   `json:"max_qps"`                // 压测停止的最大QPS
	ErrorTolerateRate    float64 `json:"error_tolerate_rate"`    // 错误容忍率
}

type ModuleConf struct {
	HpaExecutorConf NormalTesterConf `json:"hpa_executor_conf"`
}

func initModuleConf() {
	var conf ModuleConf
	var err error

	d, err := os.ReadFile(moduleConf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	err = jsoniter.Unmarshal(d, &conf)
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
	DefaultClientTimeout = conf.HpaExecutorConf.DefaultClientTimeout
	RtTestEpoch = conf.HpaExecutorConf.RtTestEpoch
	MaxQPS = conf.HpaExecutorConf.MaxQPS
	ErrorTolerateRate = conf.HpaExecutorConf.ErrorTolerateRate
}
