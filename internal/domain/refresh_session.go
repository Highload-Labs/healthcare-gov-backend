package domain

import (
	"database/sql"
	"time"
)

type RefreshSession struct {
	ID        string       `json:"id" db:"id"`
	UserID    string       `json:"user_id" db:"user_id"`
	TokenHash string       `json:"token_hash" db:"token_hash"`
	ExpiresAt time.Time    `json:"expires_at" db:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at" db:"revoked_at"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}
