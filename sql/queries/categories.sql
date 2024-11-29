-- CATEGORY
-- name: CreateCategory :one
INSERT INTO categories (
		name, 
		parent_id
		root_id, 
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at;

-- name: GetCategory :one
SELECT 
		id,
		name, 
		parent_id
		root_id, 
		created_at,
		modified_at,
	FROM categories
	WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1



-- Movie_category & Movie_keywords
-- name: CreateMovie_category :one
INSERT INTO movie_categories (
		name, 
		parent_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at

-- name: CreateMovie_keywords :one
INSERT INTO movie_keywords (
		name, 
		parent_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at

-- name: GetMovie_category :one
SELECT 
		name, 
		parent_id
	FROM movie_category
	WHERE id = $1;

-- name: GetMovie_keywords :one
SELECT 
		name, 
		parent_id
	FROM Movie_keywords
	WHERE id = $1;

-- name: DeleteMovie_category :exec
DELETE FROM categories WHERE id = $1

-- name: DeleteMovie_keywords :exec
DELETE FROM categories WHERE id = $1

