package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryMiddleware(t *testing.T) {
	m := RecoveryMiddleware("development")

	handler := m(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				panic("intended panic")
			},
		),
	)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}
