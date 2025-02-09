package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RequestCounter counts total HTTP requests, labeled by path, method, and status.
var RequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path", "method", "status"},
)

// RequestDuration measures the duration (in seconds) of HTTP requests.
var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of response latency for HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method", "status"},
)

func init() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
}

// InstrumentationMiddleware records metrics for each HTTP request.
func InstrumentationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Wrap the ResponseWriter to capture the status code.
		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		status := http.StatusText(rw.status)

		RequestCounter.WithLabelValues(path, method, status).Inc()
		RequestDuration.WithLabelValues(path, method, status).Observe(duration)
	})
}

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsHandler exposes the /metrics endpoint.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
