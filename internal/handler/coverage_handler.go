package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

func (h *Handler) CoverageGetByZipcodeHandler(w http.ResponseWriter, r *http.Request) {
	zipcode := r.PathValue("zipcode")

	var req dto.CoverageRequest
	req.Zipcode = zipcode

	err := req.Validate()
	if err != nil {
		shared.SendJSONError(w, shared.ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	coverage, err := h.coverageService.GetCoverageByZipcode(r.Context(), zipcode)
	if err != nil {
		if errors.Is(err, service.ErrCoverageNotFound) {
			shared.SendJSONError(
				w,
				shared.ErrorResponse{Message: "Coverage is not available for this zipcode."},
				http.StatusNotFound,
			)
			return
		}

		shared.SendJSONError(w, shared.ErrorResponse{Message: "Internal Server Error."}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		&dto.CoverageResponse{
			Success: true,
			Data: &dto.CoverageResponseData{
				Zipcode:   zipcode,
				State:     coverage.State,
				Supported: true,
			},
		},
	)
}
