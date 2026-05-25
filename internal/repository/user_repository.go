package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (string, error)
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserRepositoryImpl struct {
	postgres *infra.Postgresql
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user domain.User) (string, error) {
	err := r.postgres.Db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id",
		user.Email,
		user.Username,
		user.Password,
	).Scan(&user.ID)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	user := domain.User{}

	row := r.postgres.Db.QueryRowContext(
		ctx,
		"SELECT id, email, username, password  FROM users WHERE id = $1",
		userID,
	)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := domain.User{}

	row := r.postgres.Db.QueryRowContext(
		ctx,
		"SELECT id, email, username, password  FROM users WHERE email = $1",
		email,
	)
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func NewUserRepository(postgres *infra.Postgresql) UserRepository {
	return &UserRepositoryImpl{
		postgres: postgres,
	}
}
