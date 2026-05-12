package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestId, _ := r.Context().Value(RequestIDKey).(string)

			next.ServeHTTP(w, r)

			log.Printf(
				`request_id=%s method=%s path=%s remote_addr=%s duration=%s`,
				requestId, r.Method, r.URL.Path, r.RemoteAddr, time.Since(start),
			)
		},
	)
}
