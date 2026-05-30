CREATE TABLE IF NOT EXISTS refresh_sessions
(
    id         UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_id    UUID      NOT NULL,
    token_hash TEXT      NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)