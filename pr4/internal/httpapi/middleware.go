package httpapi

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/pz4-monitoring/internal/metrics"
)

// normalizePath заменяет конкретные ID в пути на шаблон {id},
// чтобы не плодить отдельные label-значения для каждого студента.
func normalizePath(path string) string {
	if strings.HasPrefix(path, "/students/") {
		return "/students/{id}"
	}
	return path
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Доп. задание 3: Gauge — увеличиваем в начале, уменьшаем после обработки
		metrics.ActiveRequests.Inc()
		defer metrics.ActiveRequests.Dec()

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		path := normalizePath(r.URL.Path)

		metrics.HttpRequestsTotal.WithLabelValues(r.Method, path).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		if lrw.StatusCode() >= 400 {
			metrics.HttpErrorsTotal.WithLabelValues(
				r.Method,
				path,
				strconv.Itoa(lrw.StatusCode()),
			).Inc()
		}
	})
}
