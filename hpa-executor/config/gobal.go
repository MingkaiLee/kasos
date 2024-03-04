package config

import "sync"

var (
	// 全局的deployment-replica数量缓存表
	deploymentReplicaCache map[string]int32
	// 读写锁
	muDR sync.RWMutex
	// 全局的service-临界QPS缓存表
	serviceThreshQPSCache map[string]int
	// 读写锁
	muSTQ sync.RWMutex
)

func initGlobal() {
	deploymentReplicaCache = make(map[string]int32)
	serviceThreshQPSCache = make(map[string]int)
}

func UpdateDeploymentReplicaCache(serviceName string, replica int32) {
	muDR.Lock()
	deploymentReplicaCache[serviceName] = replica
	muDR.Unlock()
}

func UpDateServiceThreshQPSCache(serviceName string, threshQPS int) {
	muSTQ.Lock()
	serviceThreshQPSCache[serviceName] = threshQPS
	muSTQ.Unlock()
}

func GetDeplotmentReplica(serviceName string) (int32, bool) {
	muDR.RLock()
	r, ok := deploymentReplicaCache[serviceName]
	muDR.RUnlock()
	return r, ok
}

func GetServiceThreshQPS(serviceName string) (int, bool) {
	muSTQ.RLock()
	r, ok := serviceThreshQPSCache[serviceName]
	muSTQ.RUnlock()
	return r, ok
}
