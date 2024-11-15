-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres) 
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMovieById :one
SELECT id, created_at, title, year, runtime, genres, version FROM movies
WHERE id = $1;

-- name: UpdateMovie :one
UPDATE movies
SET title = $2, year = $3, runtime = $4, genres = $5, version = version + 1
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;