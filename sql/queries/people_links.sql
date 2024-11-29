
-- name: create_people_link :one
INSERT INTO people_links (
		source,  
		key,
		person_id,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at, modified_at;

-- name: Get_people_links_by_person_id :many
SELECT 
		source,  
		key,
		person_id,
		language
	FROM people_links
	WHERE person_id = $1;

-- name: Delete_people_link :exec
DELETE FROM people_links WHERE person_id = $1 AND language = $2 AND key = $3;

