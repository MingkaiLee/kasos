package internal

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
