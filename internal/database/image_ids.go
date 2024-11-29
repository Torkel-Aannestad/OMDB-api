package database

import (
	"context"
	"database/sql"
	"time"
)

type ImageID struct {
	ID           int64     `json:"id"`
	ObjectID     int64     `json:"object_id"`
	ObjectType   string    `json:"object_type"`
	ImageVersion int32     `json:"image_version"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
}

type ImageIDsModel struct {
	DB *sql.DB
}

func (m ImageIDsModel) Insert(imageID *ImageID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO image_ids (
		object_id,  
		object_type,
		image_version
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at`

	args := []any{imageID.ObjectID, imageID.ObjectType, imageID.ImageVersion}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&imageID.ID,
		&imageID.CreatedAt,
		&imageID.ModifiedAt,
	)
}

func (m ImageIDsModel) GetImageForObject(movieID int64, objectType string) ([]*ImageID, error) {
	if movieID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		id,  
		object_id,  
		object_type,
		image_version
	FROM image_ids
	WHERE object_id = $1 AND object_type = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imageIDs := []*ImageID{}

	for rows.Next() {
		var imageID ImageID

		err := rows.Scan(
			&imageID.ID,
			&imageID.ObjectID,
			&imageID.ObjectType,
			&imageID.CreatedAt,
			&imageID.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		imageIDs = append(imageIDs, &imageID)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return imageIDs, nil
}

func (m ImageIDsModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM image_ids WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return ErrRecordNotFound
	}
	return nil
}
