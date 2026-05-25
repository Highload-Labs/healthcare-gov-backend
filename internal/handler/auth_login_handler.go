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
			shared.ErrorResponse{
				Success: false,
				Message: "Missing required fields",
			},
			http.StatusBadRequest,
		)
		return
	}

	if err = req.Validate(); err != nil {
		shared.SendJSONError(
			w,
			shared.ErrorResponse{
				Success: false,
				Message: "Missing required fields",
			},
			http.StatusBadRequest,
		)
		return
	}

	accessToken, refreshToken, err := h.authService.Login(
		r.Context(), service.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			shared.SendJSONError(
				w,
				shared.ErrorResponse{
					Success: false,
					Message: "Incorrect email or password.",
				},
				http.StatusUnauthorized,
			)
			return
		}

		shared.SendJSONError(
			w,
			shared.ErrorResponse{
				Success: false,
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		panic(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(
		&dto.AuthLoginResponse{
			Success: true,
			Data: &dto.AuthLoginResponseData{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresIn:    int64(h.config.AccessTokenExpired.Seconds()),
			},
		},
	)

	if err != nil {
		shared.SendJSONError(
			w,
			shared.ErrorResponse{
				Success: false,
				Message: "Internal Server Error.",
			},
			http.StatusInternalServerError,
		)
		panic(err)
		return
	}
}
