package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(m *infra.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				rw := &responseWriter{w, http.StatusOK}

				start := time.Now()
				next.ServeHTTP(rw, r)
				duration := time.Since(start).Seconds()

				pattern := r.Pattern
				pattern = strings.TrimPrefix(pattern, r.Method+" ")
				if r.Pattern == "" {
					pattern = r.URL.Path
				}

				if rw.status >= 500 {
					m.HttpRequestsFailedTotal.WithLabelValues(r.Method, pattern).Inc()
				}

				m.HttpRequestsTotal.WithLabelValues(strconv.Itoa(rw.status), r.Method, pattern).Inc()
				m.HttpRequestsDuration.WithLabelValues(r.Method, pattern).Observe(duration)
			},
		)
	}
}
