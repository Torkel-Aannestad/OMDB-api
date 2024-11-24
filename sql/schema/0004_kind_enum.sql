-- +goose Up
CREATE TYPE kind AS ENUM (
'movie',
'series',
'season',
'episode',
'movieseries'
);

-- +goose Down
DROP TYPE kind;

