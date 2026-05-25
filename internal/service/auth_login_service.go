package service

import (
	"context"
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginService interface {
	Login(ctx context.Context, input LoginInput) (accessToken string, refreshToken string, err error)
}

type LoginInput struct {
	Email    string
	Password string
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *AuthServiceImpl) Login(ctx context.Context, input LoginInput) (
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

	accessToken, err = s.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return "", "", err
	}

	expiresRefresh := time.Now().Add(s.config.RefreshTokenExpired)
	refreshToken, err = s.GenerateRefreshToken(user.ID, expiresRefresh)
	if err != nil {
		return "", "", err
	}

	hashedRefreshToken, err := shared.Hash(refreshToken)
	if err != nil {
		return
	}

	err = s.refreshSessionRepository.Create(ctx, user.ID, hashedRefreshToken, expiresRefresh)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
