package service

import (
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthJwtService interface {
	GenerateAccessToken(userId, email, username string) (string, error)
	GenerateRefreshToken(userId string, expiresAt time.Time) (string, error)
	VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error)
}

var ErrInvalidToken = errors.New("invalid token")

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
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.config.JwtAccessSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthServiceImpl) GenerateRefreshToken(userId string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    shared.ISSUER,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	)

	tokenString, err := token.SignedString(s.config.JwtRefreshSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthServiceImpl) VerifyRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return s.config.JwtRefreshSigningKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
