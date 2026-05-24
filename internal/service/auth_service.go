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

	userRepository repository.UserRepository
}

func NewAuthService(config *config.Config, userRepository repository.UserRepository) AuthService {
	return &AuthServiceImpl{config: config, userRepository: userRepository}
}
