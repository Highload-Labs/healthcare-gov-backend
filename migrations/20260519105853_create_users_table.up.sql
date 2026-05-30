CREATE TABLE IF NOT EXISTS users
(
    id         UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    email      VARCHAR(255) UNIQUE NOT NULL,
    username   VARCHAR(20)         NOT NULL,
    password   TEXT                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)