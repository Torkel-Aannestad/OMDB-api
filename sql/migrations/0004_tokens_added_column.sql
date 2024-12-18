-- +goose Up
ALTER TABLE tokens ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE tokens ADD COLUMN data JSONB not null default '{}'::jsonb;

-- +goose Down
ALTER TABLE IF EXISTS tokens DROP created_at;
ALTER TABLE IF EXISTS tokens DROP data;

