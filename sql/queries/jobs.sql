
-- name: insert_job :one
INSERT INTO jobs (
		name
	)
	VALUES ($1)
	RETURNING id, created_at, modified_at;

-- name: get_job_by_id :one
SELECT 
		id,  
		name,
		created_at,
		modified_at
	FROM jobs
	WHERE id = $1;

-- name: Delete_job :exec
DELETE FROM jobs WHERE id = $1 ;

