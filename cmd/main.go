package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
)

func main() {
	pg := infra.NewPostgresql(config.GetConfig())

	// userRepository := memory.NewUserRepository()
	userRepository := repository.NewUserRepository(pg)
	refreshSessionRepository := repository.NewRefreshTokenRepository(pg)

	authService := service.NewAuthService(config.GetConfig(), userRepository, refreshSessionRepository)

	server := transport.NewHTTP(authService)
	server.Serve()
}
