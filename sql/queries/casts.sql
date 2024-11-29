
-- name: create_cast :one
INSERT INTO casts (
		movie_id,
		person_id,
		job_id,
		role,
		position
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING created_at, modified_at;


-- name: get_casts_by_movie_id :many
SELECT 
		movie_id,
		person_id,
		job_id,
		role,
		position
	FROM casts
	WHERE movie_id = $1;

-- name: get_casts_by_person_id :many
SELECT 
		movie_id,
		person_id,
		job_id,
		role,
		position
	FROM casts
	WHERE person_id = $1;


-- name: Delete_cast :exec
DELETE FROM casts WHERE movie_id = $1 AND person_id = $2 AND job_id = $3;


