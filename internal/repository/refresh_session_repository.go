package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Highload-Labs/healthcare-gov-backend/internal/domain"
	"github.com/Highload-Labs/healthcare-gov-backend/internal/infra"
)

var ErrSessionNotFound = errors.New("session not found")

type RefreshSessionRepository interface {
	Create(ctx context.Context, userID, refreshToken string, expiresAt time.Time) error
	Revoke(ctx context.Context, tokenHash string) error
	FindByHash(ctx context.Context, tokenHash string) (*domain.RefreshSession, error)
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

func (r *RefreshSessionRepositoryImpl) Revoke(
	ctx context.Context,
	tokenHash string,
) error {
	res, err := r.postgres.Db.ExecContext(
		ctx,
		"UPDATE refresh_sessions SET revoked_at = now() WHERE token_hash = $1",
		tokenHash,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (r *RefreshSessionRepositoryImpl) FindByHash(ctx context.Context, tokenHash string) (
	*domain.RefreshSession,
	error,
) {
	session := &domain.RefreshSession{}
	err := r.postgres.Db.QueryRowContext(
		ctx,
		"SELECT user_id, token_hash, revoked_at FROM refresh_sessions WHERE token_hash = $1 AND revoked_at IS NULL",
		tokenHash,
	).Scan(&session.UserID, &session.TokenHash, &session.RevokedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	return session, nil
}
