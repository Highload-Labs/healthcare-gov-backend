package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestIDMiddleware(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				id := r.Context().Value(RequestIDKey)
				if id == nil {
					t.Fatal("request id not set")
				}
			},
		),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
}

func BenchmarkRequestIDMiddleware(b *testing.B) {
	handler := RequestIDMiddleware(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	for b.Loop() {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}
