package handler

import (
	"net/http"
)

type Handler struct {
	mux *http.ServeMux
}

func NewHandler(mux *http.ServeMux) *Handler {
	return &Handler{mux: mux}
}

func (h *Handler) InitializeRoutes() {
}
