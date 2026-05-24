package dto

import (
	"errors"
	"strings"
)

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthLoginResponse struct {
	Success bool                   `json:"success"`
	Data    *AuthLoginResponseData `json:"data,omitempty"`
}

func (r *AuthLoginRequest) Validate() error {
	if r.Email == "" || !strings.Contains(r.Email, "@") || len(r.Email) > 255 {
		return errors.New("invalid email")
	}

	if len(r.Password) < 8 || len(r.Password) >= 72 {
		return errors.New("invalid password. must be between 8-72 characters")
	}

	return nil
}
