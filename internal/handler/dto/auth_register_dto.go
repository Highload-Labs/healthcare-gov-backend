package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
)

type AuthRegisterRequest struct {
	Email    string `json:"email,required"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRegisterResponse struct {
	Success bool                      `json:"success"`
	Data    *AuthRegisterResponseData `json:"data,omitempty"`
	Code    string                    `json:"code,omitempty"`
	Message string                    `json:"message,omitempty"`
}

type AuthRegisterResponseData struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *AuthRegisterRequest) Validate() error {
	if r.Email == "" || !strings.Contains(r.Email, "@") || len(r.Email) > 255 {
		return errors.New("invalid email")
	}

	if r.Username == "" || len(r.Username) <= 3 {
		return errors.New("invalid username")
	}

	if len(r.Password) < 8 || len(r.Password) >= 72 {
		return errors.New("invalid password. must be between 8-72 characters and an alphanumeric")
	}

	passCombination := shared.CheckAlphanumeric(r.Password)
	if !passCombination {
		return errors.New("invalid password. must be between 8-72 characters and an alphanumeric")
	}

	return nil
}
