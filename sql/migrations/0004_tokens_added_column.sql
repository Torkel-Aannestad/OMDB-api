-- +goose Up
ALTER TABLE IF EXISTS tokens ADD created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE IF EXISTS tokens DROP created_at;

