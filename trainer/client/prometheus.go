package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/MingkaiLee/kasos/server/util"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var prometheusClient api.Client

const (
	qpsQueryStep   = 15
	qpsQueryPromQL = `increase{auto_hpa="on", service_name="%s"}(service_qps[15s])`
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

func FetchSerialData(ctx context.Context, startTime, endTime time.Time, serviceName string) (data map[string][]SerialDataPoint, err error) {
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
	data = make(map[string][]SerialDataPoint)
	for _, stream := range mat {
		d := make([]SerialDataPoint, 0, len(stream.Values))
		for _, v := range stream.Values {
			d = append(d, SerialDataPoint{
				Timestamp: v.Timestamp.Time().Format(time.RFC3339),
				Value:     float64(v.Value),
			})
			data[stream.Metric.String()] = d
		}
	}
	return
}
