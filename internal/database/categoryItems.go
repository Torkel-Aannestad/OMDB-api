package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var errCategoryItemTableName = errors.New("table value must be movie_keywords or movie_categories")

type CategoryItem struct {
	MovieId    int64     `json:"movie_id"`
	CategoryId int64     `json:"category_id"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type CategoryItemsModel struct {
	DB *sql.DB
}

func categoryTableNameValidation(tableName string) error {
	if tableName != "movie_categories" && tableName != "movie_keywords" {
		return errCategoryItemTableName
	}
	return nil
}

func (m CategoryItemsModel) Insert(categoryItem *CategoryItem, tableName string) error {
	err := categoryTableNameValidation(tableName)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
	INSERT INTO %v (
		movie_id,  
		category_id
	)
	VALUES ($1, $2)
	RETURNING created_at, modified_at`, tableName)

	args := []any{
		categoryItem.MovieId,
		categoryItem.CategoryId,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&categoryItem.CreatedAt,
		&categoryItem.ModifiedAt,
	)
}

func (m CategoryItemsModel) Get(movieId int64, tableName string) ([]*CategoryItem, error) {
	if movieId < 0 {
		return nil, ErrRecordNotFound
	}
	err := categoryTableNameValidation(tableName)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf(`
		SELECT 
			movie_id,  
			category_id
		FROM %v
		WHERE movie_id = $1`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryItems := []*CategoryItem{}

	for rows.Next() {
		var categoryItem CategoryItem

		err := rows.Scan(
			&categoryItem.MovieId,
			&categoryItem.CategoryId,
			&categoryItem.CreatedAt,
			&categoryItem.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}
		categoryItems = append(categoryItems, &categoryItem)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return categoryItems, nil
}

func (m CategoryItemsModel) Delete(movieID, categoryID int64, tableName string) error {
	if movieID < 0 || categoryID < 0 {
		return ErrRecordNotFound
	}
	err := categoryTableNameValidation(tableName)
	if err != nil {
		return err
	}

	stmt := fmt.Sprintf(`DELETE FROM %v WHERE movie_id = $1 AND category_id = $2`, tableName)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, movieID, categoryID)
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
