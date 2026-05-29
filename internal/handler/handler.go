package handler

import (
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/config"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
)

type Handler struct {
	mux    *http.ServeMux
	config *config.Config

	authService     service.AuthService
	coverageService service.CoverageService
}

func NewHandler(
	mux *http.ServeMux,
	cfg *config.Config,
	authService service.AuthService,
	coverageService service.CoverageService,
) *Handler {
	return &Handler{
		mux:             mux,
		config:          cfg,
		authService:     authService,
		coverageService: coverageService,
	}
}

func (h *Handler) InitializeRoutes() {
	h.mux.HandleFunc("GET /healthz", h.HealthzGetHandler)
	h.mux.HandleFunc("POST /auth/register", h.AuthRegisterPostHandler)
	h.mux.HandleFunc("POST /auth/login", h.AuthLoginPostHandler)
	h.mux.HandleFunc("GET /auth/refresh", h.AuthRefreshHandler)

	h.mux.HandleFunc("GET /coverage/{zipcode}", h.CoverageGetByZipcodeHandler)
}
