package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/service"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/golang-jwt/jwt/v5"
)

type AuthorizationMiddleware struct {
	AuthService service.AuthService
}

func (m *AuthorizationMiddleware) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			val := r.Header.Get("Authorization")
			if val == "" {
				shared.SendJSONError(
					w, shared.ErrorResponse{
						Success: false,
						Message: "Invalid or expired access token.",
					}, http.StatusUnauthorized,
				)
				return
			}

			splitVal := strings.Split(val, "Bearer ")
			if len(splitVal) != 2 {
				shared.SendJSONError(
					w, shared.ErrorResponse{
						Success: false,
						Message: "Invalid or expired access token.",
					}, http.StatusUnauthorized,
				)
				return
			}

			tokenString := splitVal[1]
			if tokenString == "" {
				shared.SendJSONError(
					w, shared.ErrorResponse{
						Success: false,
						Message: "Invalid or expired access token.",
					}, http.StatusUnauthorized,
				)
				return
			}

			claims, err := m.AuthService.VerifyAccessToken(tokenString)
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

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
