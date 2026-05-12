package handler

import (
	"encoding/json"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzGetHandler(t *testing.T) {
	h := &Handler{}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	rr := httptest.NewRecorder()

	h.HealthzGetHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", rr.Code)
	}

	var body dto.HealthzResponse

	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !body.Status {
		t.Fatalf("Expected 'status' to be true, got false")
	}
}

func BenchmarkHealthzGetHandler(b *testing.B) {
	h := &Handler{}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	for b.Loop() {
		rr := httptest.NewRecorder()

		h.HealthzGetHandler(rr, req)
	}
}
