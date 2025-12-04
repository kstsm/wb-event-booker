-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(32) NOT NULL,
    email       VARCHAR(32) NOT NULL UNIQUE,
    telegram_id BIGINT UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users (telegram_id);

-- +goose Down
DROP TABLE IF EXISTS users;






