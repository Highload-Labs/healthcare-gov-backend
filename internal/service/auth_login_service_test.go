package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_LoginSuccess(t *testing.T) {
	repo := &mockRepo{}

	inputEmail := "test@gmail.com"
	inputPassword := "test123456"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 4)
	if err != nil {
		t.Error(err)
	}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return &domain.User{Email: email, Password: string(hashedPassword)}, nil
	}

	svc := NewAuthLoginService(repo)

	err = svc.Login(context.Background(), LoginInput{Email: inputEmail, Password: inputPassword})
	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}
}

func TestLogin_LoginFailed(t *testing.T) {
	repo := &mockRepo{}

	inputEmail := "test@gmail.com"
	inputPassword := "test123456"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), 4)
	if err != nil {
		t.Error(err)
	}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return &domain.User{Email: email, Password: string(hashedPassword)}, nil
	}

	svc := NewAuthLoginService(repo)

	err = svc.Login(context.Background(), LoginInput{Email: inputEmail, Password: "test123"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockRepo{}

	repo.findByEmailFunc = func(email string) (*domain.User, error) {
		return nil, repository.ErrUserNotFound
	}

	svc := NewAuthLoginService(repo)

	err := svc.Login(
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), 4)

	repo := &mockRepo{
		findByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{Email: email, Password: string(hashedPassword)}, nil
		},
	}

	svc := NewAuthLoginService(repo)
	input := LoginInput{Email: inputEmail, Password: inputPassword}

	for b.Loop() {
		_ = svc.Login(context.Background(), input)
	}
}
