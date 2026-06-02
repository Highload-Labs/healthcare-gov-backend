package infra

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	HttpRequestsTotal       *prometheus.CounterVec
	HttpRequestsFailedTotal *prometheus.CounterVec
	HttpRequestsDuration    *prometheus.HistogramVec
	CacheRequestsTotal      *prometheus.CounterVec
}

var metrics *Metrics
var once sync.Once

func NewMetrics(reg prometheus.Registerer) *Metrics {
	once.Do(
		func() {
			metrics = &Metrics{
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

				CacheRequestsTotal: prometheus.NewCounterVec(
					prometheus.CounterOpts{
						Name: "cache_requests_total",
						Help: "Total cache looksup split by layer identity and result status",
					},
					[]string{"cache_name", "result"},
				),
			}

			reg.MustRegister(metrics.HttpRequestsTotal)
			reg.MustRegister(metrics.HttpRequestsFailedTotal)
			reg.MustRegister(metrics.HttpRequestsDuration)
			reg.MustRegister(metrics.CacheRequestsTotal)
		},
	)
	return metrics
}
