package infra

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HttpRequestsTotal       *prometheus.CounterVec
	HttpRequestsFailedTotal *prometheus.CounterVec
	HttpRequestsDuration    *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests",
			},
			[]string{"status", "method", "path"},
		),

		HttpRequestsFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "http_requests_total_failed", Help: "Total failed HTTP requests"},
			[]string{"status", "method", "path"},
		),

		HttpRequestsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_requests_duration_seconds",
				Help:    "Http requests duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}

	reg.MustRegister(m.HttpRequestsTotal)
	reg.MustRegister(m.HttpRequestsFailedTotal)
	reg.MustRegister(m.HttpRequestsDuration)

	return m
}
