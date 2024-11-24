-- +goose Up
CREATE TABLE IF NOT EXISTS movies (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    year bigint NOT NULL,
    runtime bigint NOT NULL,
    genres text[] NOT NULL,
    version bigint NOT NULL DEFAULT 1
);

-- +goose Down
DROP TABLE IF EXISTS movies;
