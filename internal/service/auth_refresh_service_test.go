package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
)

func TestRefreshSession_VerifySuccess(t *testing.T) {
	userID := "1"

	repo := &mockUserRepo{}
	sessionRepo := &mockSessionRepo{
		createFunc: func(userID, refreshToken string, expiresAt time.Time) error {
			return nil
		},
		revokeFunc: func(userID string) error {
			return nil
		},
		findByHashFunc: func(hash string) (*domain.RefreshSession, error) {
			return &domain.RefreshSession{UserID: userID}, nil
		},
	}

	cfg := &config.Config{
		JwtAccessSigningKey:  []byte("test"),
		JwtRefreshSigningKey: []byte("test"),
		AccessTokenExpired:   3600,
		BcryptCost:           4,
	}

	repo.findByIDFunc = func(id string) (*domain.User, error) {
		return &domain.User{ID: id, Email: "test@gmail.com", Username: "test"}, nil
	}

	svc := &AuthServiceImpl{config: cfg, userRepository: repo, refreshSessionRepository: sessionRepo}
	accessToken, refreshToken, err := svc.RefreshSession(context.Background(), "test", "1")
	if err != nil {
		t.Fatal(err)
	}

	if accessToken == "" {
		t.Fatal("accessToken is empty")
	}

	if refreshToken == "" {
		t.Fatal("refreshToken is empty")
	}
}

func TestRefreshSession_VerifyFailed(t *testing.T) {
	sessionRepo := &mockSessionRepo{
		findByHashFunc: func(hash string) (*domain.RefreshSession, error) {
			return nil, sql.ErrNoRows
		},
	}

	cfg := &config.Config{
		JwtAccessSigningKey:  []byte("test"),
		JwtRefreshSigningKey: []byte("test"),
		AccessTokenExpired:   3600,
		BcryptCost:           4,
	}

	svc := &AuthServiceImpl{config: cfg, refreshSessionRepository: sessionRepo}
	accessToken, refreshToken, err := svc.RefreshSession(context.Background(), "test", "1")
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatal("expected sql.ErrNoRows error")
		}
	}

	if accessToken != "" {
		t.Fatal("expected empty access token")
	}

	if refreshToken != "" {
		t.Fatal("expected empty refresh token")
	}
}
