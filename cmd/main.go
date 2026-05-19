package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository/memory"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
)

func main() {
	userRepository := memory.NewUserRepository()

	authJwtSvc := service.NewAuthJwtService(config.GetConfig())
	authRegisterSvc := service.NewAuthRegisterService(config.GetConfig(), userRepository)
	authLoginSvc := service.NewAuthLoginService(userRepository, authJwtSvc)

	server := transport.NewHTTP(authRegisterSvc, authLoginSvc)
	server.Serve()
}
