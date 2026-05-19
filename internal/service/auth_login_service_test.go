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

type mockAuthJwtService struct {
	generateAccessTokenFunc  func(userId, email, username string) (string, error)
	generateRefreshTokenFunc func(userId string) (string, error)
}

func (m *mockAuthJwtService) GenerateAccessToken(userId, email, username string) (string, error) {
	return m.generateAccessTokenFunc(userId, email, username)
}

func (m *mockAuthJwtService) GenerateRefreshToken(userId string) (string, error) {
	return m.generateRefreshTokenFunc(userId)
}

func TestLogin_LoginSuccess(t *testing.T) {
	repo := &mockRepo{}
	jwtService := &mockAuthJwtService{
		generateAccessTokenFunc: func(userId, email, username string) (string, error) {
			return "accessToken", nil
		},
		generateRefreshTokenFunc: func(userId string) (string, error) {
			return "refreshToken", nil
		},
	}

	cfg := &config.Config{
		BcryptCost: 4,
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

	svc := NewAuthLoginService(repo, jwtService)

	accessToken, refreshToken, err := svc.Login(
		context.Background(),
		LoginInput{Email: inputEmail, Password: inputPassword},
	)
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}

	if accessToken != "accessToken" {
		t.Errorf("expected accessToken to be 'accessToken', got %s", accessToken)
	}

	if refreshToken != "refreshToken" {
		t.Errorf("expected refreshToken to be 'refreshToken', got %s", refreshToken)
	}
}

func TestLogin_LoginFailed(t *testing.T) {
	repo := &mockRepo{}
	jwtService := &mockAuthJwtService{
		generateAccessTokenFunc: func(userId, email, username string) (string, error) {
			return "accessToken", nil
		},
		generateRefreshTokenFunc: func(userId string) (string, error) {
			return "refreshToken", nil
		},
	}

	cfg := &config.Config{
		BcryptCost: 4,
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

	svc := NewAuthLoginService(repo, jwtService)

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
	jwtService := &mockAuthJwtService{
		generateAccessTokenFunc: func(userId, email, username string) (string, error) {
			return "accessToken", nil
		},
		generateRefreshTokenFunc: func(userId string) (string, error) {
			return "refreshToken", nil
		},
	}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return nil, repository.ErrUserNotFound
	}

	svc := NewAuthLoginService(repo, jwtService)

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

	jwtService := &mockAuthJwtService{
		generateAccessTokenFunc: func(userId, email, username string) (string, error) {
			return "accessToken", nil
		},
		generateRefreshTokenFunc: func(userId string) (string, error) {
			return "refreshToken", nil
		},
	}

	cfg := &config.Config{
		BcryptCost: 4,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), cfg.BcryptCost)

	repo := &mockRepo{
		findByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{Email: email, Password: string(hashedPassword)}, nil
		},
	}

	svc := NewAuthLoginService(repo, jwtService)
	input := LoginInput{Email: inputEmail, Password: inputPassword}

	for b.Loop() {
		_, _, _ = svc.Login(context.Background(), input)
	}
}
