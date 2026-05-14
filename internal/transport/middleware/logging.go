package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestId, _ := r.Context().Value(RequestIDKey).(string)

			next.ServeHTTP(w, r)

			slog.Info(
				"request processed",
				"request_id", requestId,
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"duration", time.Since(start),
			)
		},
	)
}
