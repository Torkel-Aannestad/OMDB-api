-- +goose Up
ALTER TABLE tokens ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE IF EXISTS tokens DROP created_at;

