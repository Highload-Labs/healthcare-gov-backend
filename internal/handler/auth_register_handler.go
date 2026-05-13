package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

func (h *Handler) AuthRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		shared.SendJSONError(
			w,
			dto.AuthRegisterResponse{Success: false, Code: "BAD_REQUEST", Message: "Missing required fields"},
			http.StatusBadRequest,
		)
		return
	}

	if err = req.Validate(); err != nil {
		shared.SendJSONError(
			w,
			dto.AuthRegisterResponse{Success: false, Code: "BAD_REQUEST", Message: "Missing required fields"},
			http.StatusBadRequest,
		)
		return
	}

	err = h.authRegisterSvc.Register(
		r.Context(), service.RegisterInput{
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
		},
	)

	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyUsed) {
			shared.SendJSONError(
				w,
				dto.AuthRegisterResponse{
					Success: false,
					Code:    "CONFLICT",
					Message: "This email address is already registered.",
				},
				http.StatusBadRequest,
			)
			return
		}

		shared.SendJSONError(
			w,
			dto.AuthRegisterResponse{
				Success: false,
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(
		&dto.AuthRegisterResponse{
			Success: true,
			Data: &dto.AuthRegisterResponseData{
				UserID:    "not-yet-implemented",
				CreatedAt: time.Now(),
			},
		},
	)

	if err != nil {
		shared.SendJSONError(
			w,
			dto.AuthRegisterResponse{
				Success: false,
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		return
	}
}
