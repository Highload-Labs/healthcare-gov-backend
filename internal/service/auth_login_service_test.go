package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_LoginSuccess(t *testing.T) {
	repo := &mockRepo{}

	cfg := &config.Config{
		JwtAccessSigningKey:  []byte("test"),
		JwtRefreshSigningKey: []byte("test"),
		AccessTokenExpired:   3600,
		BcryptCost:           4,
	}

	inputEmail := "test@gmail.com"
	inputPassword := "test123456"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), cfg.BcryptCost)
	if err != nil {
		t.Error(err)
	}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return &domain.User{Email: email, Password: string(hashedPassword)}, nil
	}

	svc := &AuthServiceImpl{
		config:         cfg,
		userRepository: repo,
	}

	accessToken, refreshToken, err := svc.Login(
		context.Background(),
		LoginInput{Email: inputEmail, Password: inputPassword},
	)
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}

	if accessToken == "" {
		t.Errorf("expected accessToken to not null, got %s", accessToken)
	}

	if refreshToken == "" {
		t.Errorf("expected refreshToken to not null, got %s", refreshToken)
	}
}

func TestLogin_LoginFailed(t *testing.T) {
	repo := &mockRepo{}

	cfg := &config.Config{
		JwtAccessSigningKey:  []byte("test"),
		JwtRefreshSigningKey: []byte("test"),
		AccessTokenExpired:   3600,
		BcryptCost:           4,
	}

	inputEmail := "test@gmail.com"
	inputPassword := "test123456"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), cfg.BcryptCost)
	if err != nil {
		t.Error(err)
	}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return &domain.User{Email: email, Password: string(hashedPassword)}, nil
	}

	svc := &AuthServiceImpl{
		config:         cfg,
		userRepository: repo,
	}

	_, _, err = svc.Login(
		context.Background(),
		LoginInput{Email: inputEmail, Password: "test123"},
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockRepo{}
	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return nil, repository.ErrUserNotFound
	}

	svc := &AuthServiceImpl{
		userRepository: repo,
	}

	_, _, err := svc.Login(
		context.Background(), LoginInput{
			Email:    "missing@test.com",
			Password: "password123",
		},
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func BenchmarkLogin(b *testing.B) {
	inputEmail := "test@gmail.com"
	inputPassword := "test123456"

	cfg := &config.Config{
		JwtAccessSigningKey:  []byte("test"),
		JwtRefreshSigningKey: []byte("test"),
		AccessTokenExpired:   3600,
		BcryptCost:           4,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), cfg.BcryptCost)

	repo := &mockRepo{
		findByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{Email: email, Password: string(hashedPassword)}, nil
		},
	}

	svc := &AuthServiceImpl{
		config:         cfg,
		userRepository: repo,
	}

	input := LoginInput{Email: inputEmail, Password: inputPassword}

	for b.Loop() {
		_, _, _ = svc.Login(context.Background(), input)
	}
}
