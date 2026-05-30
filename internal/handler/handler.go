package handler

import (
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/transport/middleware"
)

type Handler struct {
	mux    *http.ServeMux
	config *config.Config

	authService     service.AuthService
	coverageService service.CoverageService
	planService     service.PlanService

	authorizationMiddleware *middleware.AuthorizationMiddleware
}

func NewHandler(
	mux *http.ServeMux,
	cfg *config.Config,
	authService service.AuthService,
	coverageService service.CoverageService,
	planService service.PlanService,
	authorizationMiddleware *middleware.AuthorizationMiddleware,
) *Handler {
	return &Handler{
		mux:                     mux,
		config:                  cfg,
		authService:             authService,
		coverageService:         coverageService,
		planService:             planService,
		authorizationMiddleware: authorizationMiddleware,
	}
}

func (h *Handler) InitializeRoutes() {
	h.mux.HandleFunc("GET /healthz", h.HealthzGetHandler)
	h.mux.HandleFunc("POST /auth/register", h.AuthRegisterPostHandler)
	h.mux.HandleFunc("POST /auth/login", h.AuthLoginPostHandler)
	h.mux.HandleFunc("GET /auth/refresh", h.AuthRefreshHandler)

	h.mux.HandleFunc("GET /coverage/{zipcode}", h.CoverageGetByZipcodeHandler)

	h.mux.HandleFunc("GET /plans", h.authorizationMiddleware.Authorization(http.HandlerFunc(h.PlansGetByZipcode)).ServeHTTP)
}
