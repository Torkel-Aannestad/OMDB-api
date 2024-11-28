package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Category struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ParentID   NullInt64 `json:"parent_id,omitempty"`
	RootID     NullInt64 `json:"root_id,omitempty"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type CategoryItem struct {
	MovieId    int64 `json:"movie_id"`
	CategoryId int64 `json:"category_id"`
}

type CategoriesModel struct {
	DB *sql.DB
}

func (m CategoriesModel) InsertCategory(category *Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO categories (
		name, 
		parent_id
		root_id, 
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at`

	args := []any{
		category.Name,
		category.ParentID.NullInt64,
		category.RootID.NullInt64,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.ModifiedAt,
	)
}

func (m CategoriesModel) GetCategory(id int64) (*Category, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	category := Category{}

	query := `
	SELECT 
		id,
		name, 
		parent_id
		root_id, 
		created_at,
		modified_at,
	FROM categories
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID.NullInt64,
		&category.RootID.NullInt64,
		&category.CreatedAt,
		&category.ModifiedAt,
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

func (m CategoriesModel) DeleteCategory(id int64) error {
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

// func ValidateCategory(v *validator.Validator, category *Category) {
// 	v.Check(category.Name != "", "name", "must be provided")
// 	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")

// 	v.Check(category.Birthday.IsZero(), "date", "must be provided")
// 	v.Check(category.Birthday.Year() >= 1888, "date", "must be greater than year 1888")
// 	v.Check(category.Birthday.Compare(time.Now()) < 1, "date", "must not be in the future")

// 	v.Check(category.Deathday.Year() >= 1888, "date", "must be greater than year 1888")
// 	v.Check(category.Deathday.Compare(time.Now()) < 1, "date", "must not be in the future")

// 	v.Check(category.Gender != "", "gender", "must be provided")
// 	v.Check(category.Gender != "male" && category.Gender != "female" && category.Gender != "non-binary" && category.Gender != "unknown", "Gender", "must be one of the following values: male, female, non-binary, unknown")

// 	v.Check(validator.Unique(category.Aliases), "aliases", "must not contain duplicate values")

// }

func (m CategoriesModel) InsertCategoryItem(categoryItem *CategoryItem, table string) error {
	if table != "movie_categories" && table != "movie_keywords" {
		return errors.New("table value must be movie_keywords or movie_categories")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// MovieId    int64 `json:"movie_id"`
	// CategoryId int64 `json:"category_id"`
	query := fmt.Sprintf(`
	INSERT INTO %v (
		name, 
		parent_id
	)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, modified_at`)

	args := []any{
		category.Name,
		category.ParentID.NullInt64,
		category.RootID.NullInt64,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.ModifiedAt,
	)
}
