package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginService interface {
	Login(ctx context.Context, input LoginInput) error
}

type AuthLoginServiceImpl struct {
	userRepository repository.UserRepository
}

func NewAuthLoginService(userRepo repository.UserRepository) AuthLoginService {
	return &AuthLoginServiceImpl{
		userRepository: userRepo,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *AuthLoginServiceImpl) Login(ctx context.Context, input LoginInput) error {
	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidCredentials
		}

		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}

		return err
	}

	return err
}
