
-- name: create_movie_categories :one
INSERT INTO movie_categories (
		movie_id,  
		category_id
	)
	VALUES ($1, $2)
	RETURNING created_at, modified_at;

-- name: create_movie_keywords :one
INSERT INTO movie_keywords (
		movie_id,  
		category_id
	)
	VALUES ($1, $2)
	RETURNING created_at, modified_at;

-- name: Get_movie_categories_by_movie_id :many
SELECT 
		movie_id,  
		category_id
	FROM movie_categories
	WHERE movie_id = $1;

-- name: Get_movie_keywords_by_movie_id :many
SELECT 
		movie_id,  
		category_id
	FROM movie_keywords
	WHERE movie_id = $1;

-- name: Delete_movie_categories :exec
DELETE FROM movie_categories WHERE movie_id = $1 AND category_id = $2;

-- name: Delete_movie_keywords :exec
DELETE FROM movie_keywords WHERE movie_id = $1 AND category_id = $2;

