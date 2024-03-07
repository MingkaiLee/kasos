package main

import "github.com/prometheus/client_golang/prometheus"

var (
	ServiceQPSMetric *prometheus.CounterVec
	LatencyMetric    *prometheus.HistogramVec
)

func init() {
	ServiceQPSMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "service_qps",
			Help: "Total number of HTTP requests processed.",
		},
		[]string{"auto_hpa", "service_name"},
	)
	// 毫秒级接口延迟的监控
	LatencyMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "latency",
			Help:    "HTTP request latency in milliseconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service_name"},
	)
	prometheus.MustRegister(ServiceQPSMetric)
	prometheus.MustRegister(LatencyMetric)
}
