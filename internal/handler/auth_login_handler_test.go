package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
)

type mockAuthLoginService struct {
	loginFunc func(input service.LoginInput) (string, string, error)
}

func (m *mockAuthLoginService) Login(ctx context.Context, input service.LoginInput) (
	accessToken string,
	refreshToken string,
	err error,
) {
	return m.loginFunc(input)
}

func TestAuthLoginPostHandler_ValidationFailure(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil)

	body, _ := json.Marshal(dto.AuthLoginRequest{Password: "pass1234"})
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	h.AuthLoginPostHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func BenchmarkAuthLogin_ValidationSuccess(b *testing.B) {
	req := dto.AuthLoginRequest{Email: "Test@gmail.com", Password: "1a1asd1231asd"}

	for b.Loop() {
		_ = req.Validate()
	}
}

func BenchmarkAuthLogin_ValidationFailure(b *testing.B) {
	req := dto.AuthLoginRequest{Email: "Test1231.com", Password: "11111111"}

	for b.Loop() {
		_ = req.Validate()
	}
}

func BenchmarkAuthLoginPostHandler_Success(b *testing.B) {
	svc := &mockAuthLoginService{
		loginFunc: func(input service.LoginInput) (string, string, error) {
			return "access", "refresh", nil
		},
	}

	cfg := &config.Config{
		AccessTokenExpired: 3600,
	}

	h := NewHandler(nil, cfg, nil, svc)

	body, _ := json.Marshal(
		dto.AuthRegisterRequest{
			Email:    "Test@gmail.com",
			Password: "1a1asd1231asd",
		},
	)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

	for b.Loop() {
		rr := httptest.NewRecorder()

		h.AuthLoginPostHandler(rr, req)
	}
}
