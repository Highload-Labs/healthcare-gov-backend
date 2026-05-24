package service

import (
	"testing"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
)

func TestAuthJwtService_GenerateAccessToken(t *testing.T) {
	cfg := &config.Config{
		AccessTokenExpired:  1 * time.Hour,
		JwtAccessSigningKey: []byte("test_access"),
	}

	service := &AuthServiceImpl{config: cfg}

	token, err := service.GenerateAccessToken("550e8400-e29b-41d4-a716-446655440000", "test@gmail.com", "test")
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}
}

func TestAuthJwtService_GenerateRefreshToken(t *testing.T) {
	cfg := &config.Config{
		AccessTokenExpired:   1 * time.Hour,
		RefreshTokenExpired:  168 * time.Hour,
		JwtAccessSigningKey:  []byte("test_access"),
		JwtRefreshSigningKey: []byte("test_refresh"),
	}

	service := &AuthServiceImpl{config: cfg}

	token, err := service.GenerateRefreshToken("550e8400-e29b-41d4-a716-446655440000")
	if err != nil {
		t.Error(err)
	}

	if token == "" {
		t.Error("token is empty")
	}
}

func BenchmarkAuthJwtService_GenerateAccessToken(b *testing.B) {
	cfg := &config.Config{
		AccessTokenExpired:   1 * time.Hour,
		RefreshTokenExpired:  168 * time.Hour,
		JwtAccessSigningKey:  []byte("test_access"),
		JwtRefreshSigningKey: []byte("test_refresh"),
	}

	service := &AuthServiceImpl{config: cfg}

	for b.Loop() {
		_, _ = service.GenerateAccessToken("a-b-c-d", "test@gmail.com", "test")
	}
}
