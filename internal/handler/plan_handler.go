package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

func (h *Handler) PlansGetByZipcode(w http.ResponseWriter, r *http.Request) {
	zipcode := r.URL.Query().Get("zipcode")

	var req dto.PlansZipcodeRequest
	req.Zipcode = zipcode

	err := req.Validate()
	if err != nil {
		shared.SendJSONError(w, shared.ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	plans, err := h.planService.GetByZipcode(r.Context(), zipcode)
	if err != nil {
		if errors.Is(err, service.ErrPlanNotFound) || errors.Is(err, service.ErrCoverageNotFound) {
			shared.SendJSONError(w, shared.ErrorResponse{Message: "No plans available for this zipcode."}, http.StatusNotFound)
			return
		}

		shared.SendJSONError(w, shared.ErrorResponse{Message: "Internal Server Error."}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.PlansResponse{
		Success: true,
		Data:    plans,
	})
}
