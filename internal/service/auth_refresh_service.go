package service

import (
	"context"
)

type AuthRefreshService interface {
	RefreshSession(ctx context.Context, userID string) (
		accessToken string,
		refreshToken string,
		err error,
	)
}

func (s *AuthServiceImpl) RefreshSession(ctx context.Context, userID string) (
	accessToken string,
	refreshToken string,
	err error,
) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return
	}

	accessToken, err = s.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		return
	}

	refreshToken, err = s.GenerateRefreshToken(user.ID)
	if err != nil {
		return
	}

	return
}
