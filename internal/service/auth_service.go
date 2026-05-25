package service

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
)

type AuthService interface {
	AuthJwtService
	AuthLoginService
	AuthRegisterService
	AuthRefreshService
}

type AuthServiceImpl struct {
	config *config.Config

	userRepository           repository.UserRepository
	refreshSessionRepository repository.RefreshSessionRepository
}

func NewAuthService(
	config *config.Config,
	userRepository repository.UserRepository,
	sessionRepository repository.RefreshSessionRepository,
) AuthService {
	return &AuthServiceImpl{config: config, userRepository: userRepository, refreshSessionRepository: sessionRepository}
}
