package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
)

func (h *Handler) HealthzGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(dto.HealthzResponse{Status: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
