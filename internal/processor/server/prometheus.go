package processorserver

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests received",
}, []string{"status", "path", "method"})

// Helper to capture HTTP status code
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Middleware to count HTTP requests
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Wrap the ResponseWriter to capture the status code
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// process the next request
		next.ServeHTTP(recorder, r)

		method := r.Method
		path := r.URL.Path // can be adjusted for specific routes
		status := strconv.Itoa(recorder.statusCode)

		httpRequestCounter.WithLabelValues(status, path, method).Inc()

	})
}

func initPrometheus(r *chi.Mux) {

	reg := prometheus.NewRegistry()
	reg.MustRegister(httpRequestCounter)

	pHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	r.Handle("/metrics", pHandler)

}
