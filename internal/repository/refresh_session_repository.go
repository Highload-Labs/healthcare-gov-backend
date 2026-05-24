package repository

import (
	"context"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

type RefreshSessionRepository interface {
	Create(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error
}

type RefreshSessionRepositoryImpl struct {
	postgres *infra.Postgresql
}

func NewRefreshTokenRepository(postgres *infra.Postgresql) RefreshSessionRepository {
	return &RefreshSessionRepositoryImpl{postgres: postgres}
}

func (r *RefreshSessionRepositoryImpl) Create(
	ctx context.Context,
	userID, tokenHash string,
	expiresAt time.Time,
) error {
	_, err := r.postgres.Db.ExecContext(
		ctx,
		"INSERT INTO refresh_sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		userID,
		tokenHash,
		expiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}
