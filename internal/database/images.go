package database

import (
	"context"
	"database/sql"
	"time"
)

type Image struct {
	ID           int64     `json:"id"`
	ObjectID     int64     `json:"object_id"`
	ObjectType   string    `json:"object_type"`
	ImageVersion int32     `json:"image_version"`
	CreatedAt    time.Time `json:"created_at"`
	ModifiedAt   time.Time `json:"modified_at"`
}

type ImagesModel struct {
	DB *sql.DB
}

func (m ImagesModel) Insert(image *Image) error {
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

	args := []any{image.ObjectID, image.ObjectType, image.ImageVersion}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&image.ID,
		&image.CreatedAt,
		&image.ModifiedAt,
	)
}

func (m ImagesModel) GetImageForObject(movieID int64, objectType string) ([]*Image, error) {
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

	images := []*Image{}

	for rows.Next() {
		var image Image

		err := rows.Scan(
			&image.ID,
			&image.ObjectID,
			&image.ObjectType,
			&image.CreatedAt,
			&image.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, &image)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (m ImagesModel) Delete(id int64) error {
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
