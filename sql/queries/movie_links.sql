
-- name: insert_movie_link :one
INSERT INTO movie_links (
		source,  
		key,
		movie_id,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at, modified_at;

-- name: get_movie_links_by_movie_id :many
SELECT 
		source,  
		key,
		movie_id,
		language
	FROM movie_links
	WHERE movie_id = $1;

-- name: Delete_movie_link :exec
DELETE FROM movie_links WHERE movie_id = $1 AND language = $2 AND key = $3;

