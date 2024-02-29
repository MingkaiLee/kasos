package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/trainer/util"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var prometheusClient api.Client

const (
	qpsQueryStep   = 15
	qpsQueryPromQL = `increase(service_qps{auto_hpa="on",service_name="%s"}[15s])`
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

func FetchSerialData(ctx context.Context, startTime, endTime time.Time, serviceName string) (data []SerialDataPoint, err error) {
	v1api := v1.NewAPI(prometheusClient)
	rng := v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  qpsQueryStep * time.Second,
	}
	query := fmt.Sprintf(qpsQueryPromQL, serviceName)
	result, warnings, err := v1api.QueryRange(ctx, query, rng)
	if err != nil {
		util.LogErrorf("range: %+v, error: %v", rng, err)
		return
	}
	if len(warnings) > 0 {
		util.LogInfof("warnings: %+v", warnings)
	}
	mat, ok := result.(model.Matrix)
	if !ok {
		util.LogErrorf("result is not matrix, result type: %s", result.Type().String())
		return
	}
	if len(mat) == 0 {
		util.LogErrorf("result matrix is empty")
		return
	}

	data = make([]SerialDataPoint, len(mat[0].Values))

	for _, v := range mat[0].Values {
		data = append(data, SerialDataPoint{
			Timestamp: v.Timestamp.Time().Format(time.DateTime),
			Value:     float64(v.Value),
		})
	}

	return
}
