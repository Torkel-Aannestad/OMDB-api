package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
	"github.com/lib/pq"
	pg "github.com/lib/pq"
)

type Movie struct {
	ID          int64     `json:"id"`
	ParentID    NullInt64 `json:"parent_id,omitempty"`
	SeriesID    NullInt64 `json:"series_id,omitempty"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date,omitempty"`
	Kind        string    `json:"kind"`
	Runtime     int64     `json:"runtime,omitempty"`
	Budget      float64   `json:"budget,omitempty"`
	Revenue     float64   `json:"revenue,omitempty"`
	Homepage    string    `json:"homepage,omitempty"`
	VoteAvarage float64   `json:"vote_average,omitempty"`
	VoteCount   int64     `json:"vote_count,omitempty"`
	Abstract    string    `json:"abstract,omitempty"`
	Version     int32     `json:"version"`
	CreatedAt   time.Time `json:"-"`
	ModifiedAt  time.Time `json:"-"`
}

type MovieModel struct {
	DB *sql.DB
}

// UPDATE movies
// SET name = $3, parent_id = $3, date = $4, series_id = $5, kind = $6, runtime = $7, budget = $8, revenue = $9, homepage = $10, vote_average = $11, votes_count = $12, abstract = $13, modified_at = $14, version = version + 1
// WHERE id = $1 and version = $2
// RETURNING *;

// DELETE FROM movies
// WHERE id = $1;

// -- SELECT count(*) OVER(), id, name, parent_id, date, series_id, kind, runtime, budget, revenue, homepage, vote_average, votes_count, abstract, created_at, modified_at, version
// -- FROM movies
// -- WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', sqlc.arg(name)) OR sqlc.arg(name) = '')
// -- AND (genres @> sqlc.arg(genres) OR sqlc.arg(genres) = '{}')
// -- ORDER BY title ASC, id ASC
// -- LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

func (m MovieModel) Insert(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO movies (
	name, 
	parent_id, 
	date, 
	series_id, 
	kind, 
	runtime, 
	budget, 
	revenue, 
	homepage, 
	vote_average, 
	votes_count, 
	abstract )
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id, created_at, modified_at, version`

	args := []any{
		movie.Name,
		movie.ParentID,
		movie.Date,
		movie.SeriesID,
		movie.Kind,
		movie.Runtime,
		movie.Budget,
		movie.Revenue,
		movie.Homepage,
		movie.VoteAvarage,
		movie.VoteCount,
		movie.Abstract,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.ModifiedAt,
		&movie.Version,
	)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	movie := Movie{}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {

	sortColumn := filters.getSortColumn()
	sortDirection := filters.getSortDirection()

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, sortColumn, sortDirection)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pg.Array(genres), filters.limit(), filters.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pg.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(filters.Page, filters.PageSize, totalRecords)

	return movies, metadata, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
	UPDATE movies 
	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		} else {
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM movies WHERE id = $1
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

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
