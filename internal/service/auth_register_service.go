package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthRegisterService interface {
	Register(ctx context.Context, input RegisterInput) error
}

type AuthRegisterServiceImpl struct {
	userRepository repository.UserRepository
}

func NewAuthRegisterService(userRepo repository.UserRepository) AuthRegisterService {
	return &AuthRegisterServiceImpl{
		userRepository: userRepo,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	Username string
}

var ErrEmailAlreadyUsed = errors.New("email already used")

func (s *AuthRegisterServiceImpl) Register(ctx context.Context, input RegisterInput) error {
	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	if user != nil {
		return ErrEmailAlreadyUsed
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), config.GetConfig().BcryptCost)
	if err != nil {
		return err
	}

	err = s.userRepository.Create(
		ctx, domain.User{
			Email:    input.Email,
			Username: input.Username,
			Password: string(hashedPassword),
		},
	)

	return err
}
