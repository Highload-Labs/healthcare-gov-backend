package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	handler := LoggingMiddleware(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal("unexpected status")
	}
}

func BenchmarkLoggingMiddleware(b *testing.B) {
	handler := LoggingMiddleware(
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
