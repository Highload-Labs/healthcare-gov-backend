package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/handler/dto"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/repository"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/golang-jwt/jwt/v5"
)

func (h *Handler) AuthRefreshHandler(w http.ResponseWriter, r *http.Request) {
	val := r.Header.Get("Authorization")
	if val == "" {
		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid or expired refresh token.",
			}, http.StatusUnauthorized,
		)
		return
	}

	splitVal := strings.Split(val, "Bearer ")
	if len(splitVal) != 2 {
		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid or expired refresh token.",
			}, http.StatusUnauthorized,
		)
		return
	}

	tokenString := splitVal[1]
	if tokenString == "" {
		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid or expired refresh token.",
			}, http.StatusUnauthorized,
		)
		return
	}

	claims, err := h.authService.VerifyRefreshToken(tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenExpired) || errors.Is(
			err,
			jwt.ErrTokenSignatureInvalid,
		) {
			shared.SendJSONError(
				w, shared.ErrorResponse{
					Success: false,
					Message: "Invalid or expired refresh token.",
				}, http.StatusUnauthorized,
			)
			return
		}

		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid Server Error.",
			}, http.StatusInternalServerError,
		)
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshSession(r.Context(), tokenString, claims.Subject)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) || errors.Is(err, repository.ErrUserNotFound) {
			shared.SendJSONError(
				w,
				shared.ErrorResponse{Success: false, Message: "Invalid or expired refresh token."},
				http.StatusUnauthorized,
			)
			return
		}

		shared.SendJSONError(
			w, shared.ErrorResponse{
				Success: false,
				Message: "Invalid Server Error.",
			}, http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		&dto.AuthRefreshResponse{
			Success: true,
			Data: &dto.AuthRefreshResponseData{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresIn:    int64(h.config.AccessTokenExpired.Seconds()),
			},
		},
	)
}
