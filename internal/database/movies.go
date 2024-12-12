package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Movie struct {
	ID          int64     `json:"id"`
	ParentID    NullInt64 `json:"parent_id"`
	SeriesID    NullInt64 `json:"series_id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Kind        string    `json:"kind"`
	Runtime     int64     `json:"runtime"`
	Budget      float64   `json:"budget"`
	Revenue     float64   `json:"revenue"`
	Homepage    string    `json:"homepage"`
	VoteAvarage float64   `json:"vote_average"`
	VoteCount   int64     `json:"vote_count"`
	Abstract    string    `json:"abstract"`
	Version     int32     `json:"version"`
	CreatedAt   time.Time `json:"-"`
	ModifiedAt  time.Time `json:"-"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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
		movie.ParentID.NullInt64,
		movie.Date,
		movie.SeriesID.NullInt64,
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
	SELECT 
		id,
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
		abstract,
		created_at,
		modified_at,
		version
	FROM movies
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.Name,
		&movie.ParentID,
		&movie.Date,
		&movie.SeriesID,
		&movie.Kind,
		&movie.Runtime,
		&movie.Budget,
		&movie.Revenue,
		&movie.Homepage,
		&movie.VoteAvarage,
		&movie.VoteCount,
		&movie.Abstract,
		&movie.CreatedAt,
		&movie.ModifiedAt,
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

func (m MovieModel) GetAll(name string, kind string, filters Filters) ([]*Movie, Metadata, error) {

	sortColumn := filters.getSortColumn()
	sortDirection := filters.getSortDirection()

	var kindFilter string
	if kind != "" {
		kindFilter = "kind = $4"
	} else {
		kindFilter = "1=1"
	}

	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, parent_id, date, series_id, kind, runtime, budget, revenue, homepage, vote_average, votes_count, abstract, created_at, modified_at, version
		FROM movies
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (%s)
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
	`, kindFilter, sortColumn, sortDirection)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	args := []any{name}
	args = append(args, filters.limit(), filters.offset())
	if kind != "" {
		args = append(args, kind)
	}

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
			&movie.Name,
			&movie.ParentID,
			&movie.Date,
			&movie.SeriesID,
			&movie.Kind,
			&movie.Runtime,
			&movie.Budget,
			&movie.Revenue,
			&movie.Homepage,
			&movie.VoteAvarage,
			&movie.VoteCount,
			&movie.Abstract,
			&movie.CreatedAt,
			&movie.ModifiedAt,
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
	SET 
		name = $3,
		parent_id = $4,
		date = $5,
		series_id = $6,
		kind = $7,
		runtime = $8,
		budget = $9,
		revenue = $10,
		homepage = $11,
		vote_average = $12,
		votes_count = $13,
		abstract = $14,
		modified_at = $15,
		version = version + 1
	WHERE id = $1 and version = $2
	RETURNING version`

	args := []any{
		&movie.ID,
		&movie.Version,
		&movie.Name,
		&movie.ParentID,
		&movie.Date,
		&movie.SeriesID,
		&movie.Kind,
		&movie.Runtime,
		&movie.Budget,
		&movie.Revenue,
		&movie.Homepage,
		&movie.VoteAvarage,
		&movie.VoteCount,
		&movie.Abstract,
		&movie.ModifiedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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
	v.Check(movie.Name != "", "name", "must be provided")
	v.Check(len(movie.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(movie.Kind == "movie" || movie.Kind == "series" || movie.Kind == "season" || movie.Kind == "episode" || movie.Kind == "movieseries", "kind", "must be one of the following values: movie, series, season, episode, movieseries")

	v.Check(!movie.Date.IsZero(), "date", "must be provided")
	v.Check(movie.Date.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(movie.Date.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Budget >= 0, "budget", "must be a positive integer")
	v.Check(movie.Revenue >= 0, "revenue", "must be a positive integer")

	v.Check(len(movie.Homepage) <= 256, "homepage", "must not be longer than 256")

	v.Check(movie.VoteAvarage >= 0, "vote_average", "must be a positive integer")
	v.Check(movie.VoteAvarage <= 10, "vote_average", "must be less or equal to 10")
	v.Check(movie.VoteCount >= 0, "votes_count", "must be a positive integer")

	v.Check(len(movie.Abstract) <= 4096, "abstract", "must not be longer than 4096")

}

func ValidateKind(v *validator.Validator, kind *string) {
	v.Check(*kind == "movie" || *kind == "series" || *kind == "season" || *kind == "episode" || *kind == "movieseries", "kind", "must be one of the following values: movie, series, season, episode, movieseries")
}

// func ValidateParentID(v *validator.Validator, int64)
// v.Check(movie.ParentID.Int64 != 0, "parent_id", "must be provided")
// v.Check(movie.ParentID.Int64 > 0, "parent_id", "must be a positive integer")

// v.Check(movie.SeriesID.Int64 != 0, "series_id", "must be provided")
// v.Check(movie.SeriesID.Int64 > 0, "series_id", "must be a positive integer")
