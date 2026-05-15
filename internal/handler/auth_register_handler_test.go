package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
)

type mockAuthRegisterService struct {
	registerFunc func(input service.RegisterInput) error
}

func (m *mockAuthRegisterService) Register(
	ctx context.Context,
	input service.RegisterInput,
) error {
	return m.registerFunc(input)
}

func TestAuthRegisterPostHandler_ValidationFailure(t *testing.T) {
	h := NewHandler(nil, nil, nil)

	// Missing Email
	body, _ := json.Marshal(dto.AuthRegisterRequest{Password: "pass1234"})
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	h.AuthRegisterPostHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func BenchmarkAuthRegister_ValidationSuccess(b *testing.B) {
	req := dto.AuthRegisterRequest{Email: "Test@gmail.com", Username: "test123", Password: "1a1asd1231asd"}

	for b.Loop() {
		_ = req.Validate()
	}
}

func BenchmarkAuthRegister_ValidationFailure(b *testing.B) {
	req := dto.AuthRegisterRequest{Email: "Test@gmail.com", Username: "test123", Password: "1111111111111"}

	for b.Loop() {
		_ = req.Validate()
	}
}

func BenchmarkAuthRegisterPostHandler_ValidationFailure(b *testing.B) {
	h := NewHandler(nil, nil, nil)

	body, _ := json.Marshal(dto.AuthRegisterRequest{Password: "pass1234"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))

	for b.Loop() {
		rr := httptest.NewRecorder()

		h.AuthRegisterPostHandler(rr, req)
	}
}

func BenchmarkAuthRegisterPostHandler_Success(b *testing.B) {
	svc := &mockAuthRegisterService{
		registerFunc: func(input service.RegisterInput) error {
			return nil
		},
	}

	h := NewHandler(nil, svc, nil)

	body, _ := json.Marshal(
		dto.AuthRegisterRequest{
			Email:    "Test@gmail.com",
			Username: "test123",
			Password: "1a1asd1231asd",
		},
	)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))

	for b.Loop() {
		rr := httptest.NewRecorder()

		h.AuthRegisterPostHandler(rr, req)
	}
}
