package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	HttpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status_code"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// Доп. задание 1: счётчик запросов по student_id
	StudentRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_student_requests_total",
			Help: "Total number of requests per student",
		},
		[]string{"student_id"},
	)

	// Доп. задание 2: отдельная histogram только для /students/{id}
	StudentRequestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "app_student_request_duration_seconds",
			Help:    "Duration of /students/{id} requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Доп. задание 3: gauge активных запросов
	ActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_active_requests",
			Help: "Number of requests currently being processed",
		},
	)
)
