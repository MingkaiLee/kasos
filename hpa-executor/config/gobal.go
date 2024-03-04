package config

var (
	// 全局的deployment-replica数量缓存表
	DeploymentReplicaCache map[string]int32
	// 全局的service-临界QPS缓存表
	ServiceThreshQPSCache map[string]int
)

func initGlobal() {
	DeploymentReplicaCache = make(map[string]int32)
	ServiceThreshQPSCache = make(map[string]int)
}
