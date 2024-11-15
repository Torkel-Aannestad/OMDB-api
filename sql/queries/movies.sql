-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMovieById :one
SELECT id, created_at, title, year, runtime, genres, version FROM movies
WHERE id = $1;
