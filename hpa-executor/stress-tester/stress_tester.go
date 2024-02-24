package stresstester

type StressTestResult struct {
	ThresholdQPS int
	Rt           int64
	P99          int64
	Avg          int64
}

type StreeTester interface {
	SetConfigByJSON(json string) error
	Run() <-chan *StressTestResult
	GetResult() *StressTestResult
}
