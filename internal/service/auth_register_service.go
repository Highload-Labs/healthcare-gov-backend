package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type AuthRegisterService interface {
	Register(ctx context.Context, input RegisterInput) (string, error)
}

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

var ErrEmailAlreadyUsed = errors.New("email already used")

func (s *AuthServiceImpl) Register(ctx context.Context, input RegisterInput) (string, error) {
	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return "", err
	}

	if user != nil {
		return "", ErrEmailAlreadyUsed
	}

	newUser := domain.User{
		Email:    input.Email,
		Username: input.Username,
		Password: input.Password,
	}

	err = newUser.HashPassword(s.config.BcryptCost)
	if err != nil {
		return "", err
	}

	userID, err := s.userRepository.Create(ctx, newUser)
	if err != nil {
		return "", err
	}

	return userID, err
}
