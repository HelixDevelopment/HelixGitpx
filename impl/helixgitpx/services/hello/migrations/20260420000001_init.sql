-- +goose Up
CREATE TABLE IF NOT EXISTS hello.greetings (
    name TEXT PRIMARY KEY,
    count BIGINT NOT NULL DEFAULT 0,
    last_said_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS hello.greetings;
