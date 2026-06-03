package main

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	pg := infra.NewPostgresql(config.GetConfig())
	infra.RunDatabaseMigrations(pg)

	redisConn := infra.NewRedisClient(config.GetConfig())

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)

	metrics := infra.NewMetrics(reg)

	// userRepository := memory.NewUserRepository()
	userRepository := repository.NewUserRepository(pg)
	coverageRepository := repository.NewCoverageRepository(pg, redisConn, metrics)
	refreshSessionRepository := repository.NewRefreshTokenRepository(redisConn)
	planRepository := repository.NewPlanRepository(pg, redisConn, metrics)
	enrollmentRepository := repository.NewEnrollmentRepository(pg)

	authService := service.NewAuthService(config.GetConfig(), userRepository, refreshSessionRepository)
	coverageService := service.NewCoverageService(coverageRepository)
	planService := service.NewPlanService(planRepository, coverageService)
	enrollmentService := service.NewEnrollmentService(enrollmentRepository)

	server := transport.NewHTTP(authService, coverageService, planService, enrollmentService, reg, metrics)
	server.Serve()
}
