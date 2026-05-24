package service

import (
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/golang-jwt/jwt/v5"
)

type AuthJwtService interface {
	GenerateAccessToken(userId, email, username string) (string, error)
	GenerateRefreshToken(userId string) (string, error)
}

type Claims struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	jwt.RegisteredClaims
}

func (s *AuthServiceImpl) GenerateAccessToken(userId, email, username string) (string, error) {
	claims := Claims{
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    shared.ISSUER,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessTokenExpired)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.config.JwtAccessSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthServiceImpl) GenerateRefreshToken(userId string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    shared.ISSUER,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenExpired)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	)

	tokenString, err := token.SignedString(s.config.JwtRefreshSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
