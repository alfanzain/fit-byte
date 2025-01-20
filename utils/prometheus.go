package utils

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Define metrics as package-level variables
var (
	HttpRequestsTotal prometheus.Counter
	ActiveConnections prometheus.Gauge
	RequestDuration   prometheus.Histogram
	Goroutines        prometheus.Gauge
)

func InitPrometheusMetrics() {
	// HTTP requests counter
	HttpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "myapp_http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
	)

	// Active connections gauge
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "myapp_active_connections",
			Help: "Current number of active connections",
		},
	)

	// Request duration histogram
	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "myapp_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
		},
	)

	// Goroutines gauge
	Goroutines = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "myapp_goroutines_total",
			Help: "Current number of Goroutines",
		},
	)

	// Register metrics
	prometheus.MustRegister(HttpRequestsTotal, ActiveConnections, RequestDuration, Goroutines)
}

func SimulateMetrics() {
	go func() {
		for {
			// Update metrics
			HttpRequestsTotal.Inc()
			ActiveConnections.Set(10)
			Goroutines.Set(float64(runtime.NumGoroutine()))
			time.Sleep(2 * time.Second)
		}
	}()
}
