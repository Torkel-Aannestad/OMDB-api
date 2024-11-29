
-- name: CreateCategory :one
INSERT INTO categories (
		name, 
		parent_id,
		root_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at;

-- name: GetCategory :one
SELECT 
		id,
		name, 
		parent_id,
		root_id, 
		created_at,
		modified_at
	FROM categories
	WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;


