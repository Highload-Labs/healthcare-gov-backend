package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository/memory"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
)

func main() {
	userRepository := memory.NewUserRepository()

	authRegisterSvc := service.NewAuthRegisterService(userRepository)
	authLoginSvc := service.NewAuthLoginService(userRepository)

	server := transport.NewHTTP(authRegisterSvc, authLoginSvc)
	server.Serve()
}
