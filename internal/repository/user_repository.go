package repository

import (
	"context"
	"errors"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
