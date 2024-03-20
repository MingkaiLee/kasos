package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/infer-module/util"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var prometheusClient api.Client

const (
	qpsQueryPromQL = `irate(service_qps{auto_hpa="on",service_name="%s"}[1m])`
)

type SerialDataPoint struct {
	Timestamp string
	Value     float64
}

func InitPrometheusClient() {
	var err error
	prometheusClient, err = api.NewClient(api.Config{
		Client:  http.DefaultClient,
		Address: "http://prometheus-k8s.monitoring.svc.cluster.local:9090",
	})
	if err != nil {
		util.LogErrorf("panic: %v", err)
		panic(err)
	}
}

// 查询服务的qps实时数据点
func RealTimeData(ctx context.Context, serviceName string) (data SerialDataPoint, err error) {
	v1api := v1.NewAPI(prometheusClient)
	query := fmt.Sprintf(qpsQueryPromQL, serviceName)
	currentTime := time.Now()
	result, warnings, err := v1api.Query(ctx, query, currentTime)
	if err != nil {
		util.LogErrorf("query: %s, error: %v", query, err)
		return
	}
	if len(warnings) > 0 {
		util.LogInfof("warnings: %+v", warnings)
	}
	value, ok := result.(*model.Scalar)
	if !ok {
		util.LogErrorf("result is not scalar, result type: %s", result.Type().String())
		return
	}
	data = SerialDataPoint{
		Timestamp: value.Timestamp.Time().Format(time.DateTime),
		Value:     float64(value.Value),
	}
	return
}
