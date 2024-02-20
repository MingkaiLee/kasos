package client

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var prometheusClient api.Client

func InitPrometheusClient() {
	var err error
	prometheusClient, err = api.NewClient(api.Config{
		Client:  http.DefaultClient,
		Address: "http://prometheus-k8s.monitoring.svc.cluster.local:9090",
	})
	if err != nil {
		panic(err)
	}
}

func FetchSerialData(ctx context.Context, startTime, endTime time.Time, step time.Duration) {
	v1api := v1.NewAPI(prometheusClient)
	rng := v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  step,
	}
	result, warnings, err := v1api.QueryRange(ctx, "node_serial_number_info", rng)
}
