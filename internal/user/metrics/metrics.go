package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var once sync.Once
var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests processed, labeled by status code, method, and endpoint.",
		},
		[]string{"code", "method", "endpoint"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func InitMetrics() {
	m := prometheus.NewRegistry()
	once.Do(func() {
		m.MustRegister(HttpRequestsTotal)
		m.MustRegister(RequestDuration)
		m.MustRegister(collectors.NewGoCollector())
		m.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	})
	m.Gather()

	//m.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	//m.Gather()
}
