
-- name: create_image_id :one
INSERT INTO image_ids (
		object_id,  
		object_type,
		image_version
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at;

-- name: get_images_for_object :many
SELECT 
		id,  
		object_id,  
		object_type,
		image_version
	FROM image_ids
	WHERE object_id = $1 AND object_type = $2;

-- name: Delete_imageid :exec
DELETE FROM image_ids WHERE id = $1;

