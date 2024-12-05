package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Category struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ParentID   NullInt64 `json:"parent_id,omitempty"`
	RootID     NullInt64 `json:"root_id,omitempty"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	Version    int32     `json:"version"`
}

type CategoriesModel struct {
	DB *sql.DB
}

func (m CategoriesModel) Insert(category *Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO categories (
		name, 
		parent_id,
		root_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at, version`

	args := []any{
		category.Name,
		category.ParentID.NullInt64,
		category.RootID.NullInt64,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.ModifiedAt,
		&category.Version,
	)
}

func (m CategoriesModel) Get(id int64) (*Category, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	category := Category{}

	query := `
		SELECT 
			id,
			name, 
			parent_id,
			root_id, 
			created_at,
			modified_at,
			version
		FROM categories
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID.NullInt64,
		&category.RootID.NullInt64,
		&category.CreatedAt,
		&category.ModifiedAt,
		&category.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoriesModel) Update(category *Category) error {
	query := `
	UPDATE categories
	SET 
		name = $3,
		parent_id = $4,
		root_id = $5,
		modified_at = NOW(),
		version = version + 1
	WHERE id = $1 and version = $2
	RETURNING version`

	args := []any{
		&category.ID,
		&category.Version,
		&category.Name,
		&category.ParentID,
		&category.RootID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&category.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		} else {
			return err
		}
	}
	return nil
}

func (m CategoriesModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM categories WHERE id = $1
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

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")
}
