package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

func (h *Handler) AuthLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		shared.SendJSONError(
			w,
			dto.AuthLoginResponse{Success: false, Code: "BAD_REQUEST", Message: "Missing required fields"},
			http.StatusBadRequest,
		)
		return
	}

	if err = req.Validate(); err != nil {
		shared.SendJSONError(
			w,
			dto.AuthLoginResponse{Success: false, Code: "BAD_REQUEST", Message: "Missing required fields"},
			http.StatusBadRequest,
		)
		return
	}

	err = h.authLoginSvc.Login(
		r.Context(), service.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			shared.SendJSONError(
				w,
				dto.AuthLoginResponse{
					Success: false,
					Code:    "UNAUTHORIZED",
					Message: "Incorrect email or password.",
				},
				http.StatusUnauthorized,
			)
			return
		}

		shared.SendJSONError(
			w,
			dto.AuthLoginResponse{
				Success: false,
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(
		&dto.AuthLoginResponse{
			Success: true,
			Data: &dto.AuthLoginResponseData{
				AccessToken:  "not-yet-implemented",
				RefreshToken: "not-yet-implemented",
				ExpiresIn:    3600,
			},
		},
	)

	if err != nil {
		shared.SendJSONError(
			w,
			dto.AuthLoginResponse{
				Success: false,
				Code:    "INTERNAL_ERROR",
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		return
	}
}
