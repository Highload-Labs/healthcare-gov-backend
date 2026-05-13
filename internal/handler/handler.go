package handler

import (
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"net/http"
)

type Handler struct {
	mux *http.ServeMux

	authRegisterSvc service.AuthRegisterService
}

func NewHandler(mux *http.ServeMux, authRegisterSvc service.AuthRegisterService) *Handler {
	return &Handler{
		mux:             mux,
		authRegisterSvc: authRegisterSvc,
	}
}

func (h *Handler) InitializeRoutes() {
	h.mux.HandleFunc("GET /healthz", h.HealthzGetHandler)
	h.mux.HandleFunc("POST /auth/register", h.AuthRegisterPostHandler)
}
