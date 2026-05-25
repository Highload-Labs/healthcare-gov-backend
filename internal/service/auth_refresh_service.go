package service

import (
	"context"
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

type AuthRefreshService interface {
	RefreshSession(ctx context.Context, rToken, userID string) (
		accessToken string,
		refreshToken string,
		err error,
	)
}

func (s *AuthServiceImpl) RefreshSession(ctx context.Context, earlyRefreshToken, userID string) (
	accessToken string,
	refreshToken string,
	err error,
) {
	hashedEarlyRefreshToken, err := shared.Hash(earlyRefreshToken)
	if err != nil {
		return
	}

	session, err := s.refreshSessionRepository.FindByHash(ctx, hashedEarlyRefreshToken)
	if err != nil {
		return
	}

	if session.UserID != userID {
		return "", "", errors.New("token payload mistmatch")
	}

	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return
	}

	accessToken, err = s.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return
	}

	expiresRefresh := time.Now().Add(s.config.RefreshTokenExpired)
	refreshToken, err = s.GenerateRefreshToken(user.ID, expiresRefresh)
	if err != nil {
		return
	}

	err = s.refreshSessionRepository.Revoke(ctx, hashedEarlyRefreshToken)

	hashedRefreshToken, err := shared.Hash(refreshToken)
	if err != nil {
		return
	}

	err = s.refreshSessionRepository.Create(ctx, user.ID, hashedRefreshToken, expiresRefresh)
	if err != nil {
		return
	}

	return
}
