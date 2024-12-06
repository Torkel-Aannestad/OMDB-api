package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Image struct {
	ID         int64     `json:"id"`
	ObjectID   int64     `json:"object_id"`
	ObjectType string    `json:"object_type"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	Version    int32     `json:"version"`
}

type ImagesModel struct {
	DB *sql.DB
}

func (m ImagesModel) Insert(image *Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO images (
		object_id,  
		object_type
	)
	VALUES ($1, $2)
	RETURNING id, created_at, modified_at, version`

	args := []any{image.ObjectID, image.ObjectType}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&image.ID,
		&image.CreatedAt,
		&image.ModifiedAt,
		&image.Version,
	)
}

func (m ImagesModel) Get(id int64) (*Image, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		id,  
		object_id,  
		object_type,
		version,
		created_at,
		modified_at
	FROM images
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var image Image
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.ObjectID,
		&image.ObjectType,
		&image.Version,
		&image.CreatedAt,
		&image.ModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &image, nil
}
func (m ImagesModel) GetImagesForObject(objectID int64, objectType string) ([]*Image, error) {
	if objectID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		id,  
		object_id,  
		object_type,
		version,
		created_at,
		modified_at
	FROM images
	WHERE object_id = $1 AND object_type = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, objectID, objectType)
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
			&image.Version,
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

func (m ImagesModel) Update(image *Image) error {
	query := `
	UPDATE images
	SET 
		object_id = $3,
		object_type = $4,
		modified_at = NOW(),
		version = version + 1
	WHERE id = $1 and version = $2
	RETURNING version`

	args := []any{
		&image.ID,
		&image.Version,
		&image.ObjectID,
		&image.ObjectType,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&image.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		} else {
			return err
		}
	}
	return nil
}

func (m ImagesModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM images WHERE id = $1
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

func ValidateImage(v *validator.Validator, image *Image) {
	v.Check(image.ObjectID != 0, "object_id", "must be provided")
	v.Check(image.ObjectID >= 0, "object_id", "must be a positive number")
	v.Check(validator.PermittedValue(image.ObjectType, "Movie", "Person", "Job", "User", "Category", "Company", "Character"), "object_type", "must be of the following values: Movie, Person, Job, User, Category, Company, Character")
}
