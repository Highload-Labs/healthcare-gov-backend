package handler

import (
	"encoding/json"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"net/http"
)

func (h *Handler) HealthzGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(dto.HealthzResponse{Status: true})
}
