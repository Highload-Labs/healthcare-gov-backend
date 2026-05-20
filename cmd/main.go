package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
	_ "github.com/lib/pq"
)

func main() {
	pg := infra.NewPostgresql(config.GetConfig())

	// userRepository := memory.NewUserRepository()
	userRepository := repository.NewUserRepository(pg)

	authJwtSvc := service.NewAuthJwtService(config.GetConfig())
	authRegisterSvc := service.NewAuthRegisterService(config.GetConfig(), userRepository)
	authLoginSvc := service.NewAuthLoginService(userRepository, authJwtSvc)

	server := transport.NewHTTP(authRegisterSvc, authLoginSvc)
	server.Serve()
}
