package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type mockRepo struct {
	findByEmailFunc func(email string) (*domain.User, error)
	createFunc      func(user domain.User) (string, error)
}

func (m *mockRepo) FindByEmail(ctx context.Context, e string) (*domain.User, error) {
	return m.findByEmailFunc(e)
}
func (m *mockRepo) Create(ctx context.Context, u domain.User) (string, error) {
	return m.createFunc(u)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := &mockRepo{
		findByEmailFunc: func(email string) (*domain.User, error) {
			return &domain.User{Email: email}, nil // User found
		},
	}

	cfg := &config.Config{
		BcryptCost: 4,
	}

	svc := &AuthServiceImpl{
		config:         cfg,
		userRepository: repo,
	}

	userID, err := svc.Register(context.Background(), RegisterInput{Email: "test@test.com"})
	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Errorf("expected ErrEmailAlreadyUsed, got %v", err)
	}

	if userID != "" {
		t.Errorf("expected empty user ID, got %v", userID)
	}
}

func BenchmarkRegister(b *testing.B) {
	repo := &mockRepo{
		findByEmailFunc: func(e string) (*domain.User, error) { return nil, repository.ErrUserNotFound },
		createFunc:      func(u domain.User) (string, error) { return "", nil },
	}

	cfg := &config.Config{
		BcryptCost: 4,
	}

	svc := &AuthServiceImpl{
		config:         cfg,
		userRepository: repo,
	}

	input := RegisterInput{Email: "a@b.com", Username: "user", Password: "password123"}

	for b.Loop() {
		_, _ = svc.Register(context.Background(), input)
	}
}
