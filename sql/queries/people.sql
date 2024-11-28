-- name: GetPerson :one
SELECT 
		id,
		name, 
		birthday
		deathday, 
		gender, 
		aliases, 
		created_at,
		modified_at,
		version
	FROM people
	WHERE id = $1;