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

// -- SELECT count(*) OVER(), id, name, parent_id, date, series_id, kind, runtime, budget, revenue, homepage, vote_average, votes_count, abstract, created_at, modified_at, version
// -- FROM movies
// -- WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', sqlc.arg(name)) OR sqlc.arg(name) = '')
// -- AND (genres @> sqlc.arg(genres) OR sqlc.arg(genres) = '{}')
// -- ORDER BY title ASC, id ASC
// -- LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

// func (m CategoriesModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {

// 	sortColumn := filters.getSortColumn()
// 	sortDirection := filters.getSortDirection()

// 	query := fmt.Sprintf(`
// 		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
// 		FROM movies
// 		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
// 		AND (genres @> $2 OR $2 = '{}')
// 		ORDER BY %s %s, id ASC
// 		LIMIT $3 OFFSET $4
// 	`, sortColumn, sortDirection)

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	args := []any{title, pg.Array(genres), filters.limit(), filters.Offset()}

// 	rows, err := m.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, Metadata{}, err
// 	}
// 	defer rows.Close()

// 	totalRecords := 0
// 	movies := []*Movie{}

// 	for rows.Next() {
// 		var movie Movie

// 		err := rows.Scan(
// 			&totalRecords,
// 			&movie.ID,
// 			&movie.CreatedAt,
// 			&movie.Title,
// 			&movie.Year,
// 			&movie.Runtime,
// 			pg.Array(&movie.Genres),
// 			&movie.Version,
// 		)
// 		if err != nil {
// 			return nil, Metadata{}, err
// 		}
// 		movies = append(movies, &movie)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, Metadata{}, err
// 	}

// 	metadata := calculateMetadata(filters.Page, filters.PageSize, totalRecords)

// 	return movies, metadata, nil
// }

func (m CategoriesModel) Update(category *Category) error {
	query := `
	UPDATE people
	SET 
		name = $3,
		birthday = $4,
		deathday = $5,
		gender = $6,
		aliases = $7,
		modified_at = $8,
		version = version + 1
	WHERE id = $1 and version = $2
	RETURNING version`

	args := []any{
		&category.ID,
		&category.Version,
		&category.Name,
		&category.Birthday,
		&category.Deathday,
		&category.Gender,
		&category.Aliases,
		&category.ModifiedAt,
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
		DELETE FROM people WHERE id = $1
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

	v.Check(category.Birthday.IsZero(), "date", "must be provided")
	v.Check(category.Birthday.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(category.Birthday.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(category.Deathday.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(category.Deathday.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(category.Gender != "", "gender", "must be provided")
	v.Check(category.Gender != "male" && category.Gender != "female" && category.Gender != "non-binary" && category.Gender != "unknown", "Gender", "must be one of the following values: male, female, non-binary, unknown")

	v.Check(validator.Unique(category.Aliases), "aliases", "must not contain duplicate values")

}
