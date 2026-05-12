package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const RequestIDKey string = "request_id"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqID := uuid.New().String()

			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
