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
	coverageRepository := repository.NewCoverageRepository(pg)
	refreshSessionRepository := repository.NewRefreshTokenRepository(pg)
	planRepository := repository.NewPlanRepository(pg)
	enrollmentRepository := repository.NewEnrollmentRepository(pg)

	authService := service.NewAuthService(config.GetConfig(), userRepository, refreshSessionRepository)
	coverageService := service.NewCoverageService(coverageRepository)
	planService := service.NewPlanService(planRepository, coverageService)
	enrollmentService := service.NewEnrollmentService(enrollmentRepository)

	server := transport.NewHTTP(authService, coverageService, planService, enrollmentService)
	server.Serve()
}
