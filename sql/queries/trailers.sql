-- name: create_trailer :one
INSERT INTO trailers (
		movie_id,
		source,  
		key,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING trailer_id, created_at, modified_at;

-- name: get_trailers_by_movie_id :many
SELECT 
		trailer_id,
		source,  
		key,
		movie_id,
		language
	FROM trailers
	WHERE movie_id = $1;

-- name: Delete_trailer :exec
DELETE FROM trailers WHERE trailer_id = $1;

