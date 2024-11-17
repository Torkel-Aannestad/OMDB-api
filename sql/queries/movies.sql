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
WHERE id = $1 and version = $6
RETURNING *;

-- name: DeleteMovie :execrows
DELETE FROM movies
WHERE id = $1;

-- name: ListMovies :many
SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', sqlc.arg(title)) OR sqlc.arg(title) = '')
AND (genres @> sqlc.arg(genres) OR sqlc.arg(genres) = '{}')
ORDER BY title ASC, id ASC
LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);
