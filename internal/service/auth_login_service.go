package service

import (
	"context"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginService interface {
	Login(ctx context.Context, input LoginInput) (accessToken string, refreshToken string, err error)
}

type AuthLoginServiceImpl struct {
	userRepository repository.UserRepository
	jwtService     AuthJwtService
}

func NewAuthLoginService(userRepo repository.UserRepository, jwtService AuthJwtService) AuthLoginService {
	return &AuthLoginServiceImpl{
		userRepository: userRepo,
		jwtService:     jwtService,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *AuthLoginServiceImpl) Login(ctx context.Context, input LoginInput) (
	accessToken string,
	refreshToken string,
	err error,
) {
	user, err := s.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", ErrInvalidCredentials
		}

		return "", "", err
	}

	if err = user.VerifyPassword(input.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", ErrInvalidCredentials
		}

		return "", "", err
	}

	accessToken, err = s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
