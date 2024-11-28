package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Person struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Birthday   time.Time `json:"birthday,omitempty"`
	Deathday   time.Time `json:"deathday,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Aliases    []string  `json:"aliases,omitempty"`
	Version    int32     `json:"version"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type PeopleModel struct {
	DB *sql.DB
}

func (m PeopleModel) Insert(person *Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO movies (
	name, 
	birthday
	deathday, 
	gender, 
	aliases, 
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, modified_at, version`

	args := []any{
		person.Name,
		person.Birthday,
		person.Deathday,
		person.Gender,
		person.Aliases,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&person.ID,
		&person.CreatedAt,
		&person.ModifiedAt,
		&person.Version,
	)
}

func (m PeopleModel) Get(id int64) (*Person, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	person := Person{}

	query := `
	SELECT 
		id,
		name, 
		birthday
		deathday, 
		gender, 
		aliases, 
		created_at,
		modified_at,
		version
	FROM people
	WHERE id = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Birthday,
		&person.Deathday,
		&person.Gender,
		&person.Aliases,
		&person.CreatedAt,
		&person.ModifiedAt,
		&person.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &person, nil
}

// -- SELECT count(*) OVER(), id, name, parent_id, date, series_id, kind, runtime, budget, revenue, homepage, vote_average, votes_count, abstract, created_at, modified_at, version
// -- FROM movies
// -- WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', sqlc.arg(name)) OR sqlc.arg(name) = '')
// -- AND (genres @> sqlc.arg(genres) OR sqlc.arg(genres) = '{}')
// -- ORDER BY title ASC, id ASC
// -- LIMIT sqlc.arg(limit_value) OFFSET sqlc.arg(offset_value);

// func (m PeopleModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {

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

func (m PeopleModel) Update(movie *Movie) error {
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

func (m PeopleModel) Delete(id int64) error {
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

func ValidatePeople(v *validator.Validator, movie *Movie) {
	v.Check(movie.Name != "", "name", "must be provided")
	v.Check(len(movie.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(movie.ParentID.Int64 != 0, "parent_id", "must be provided")
	v.Check(movie.ParentID.Int64 > 0, "parent_id", "must be a positive integer")

	v.Check(movie.SeriesID.Int64 != 0, "series_id", "must be provided")
	v.Check(movie.ParentID.Int64 > 0, "series_id", "must be a positive integer")

	v.Check(movie.Date.IsZero(), "date", "must be provided")
	v.Check(movie.Date.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(movie.Date.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Budget >= 0, "budget", "must be a positive integer")
	v.Check(movie.Revenue >= 0, "revenue", "must be a positive integer")

	v.Check(len(movie.Homepage) > 256, "homepage", "must not be longer than 256")

	v.Check(movie.VoteAvarage >= 0, "vote_average", "must be a positive integer")
	v.Check(movie.VoteAvarage <= 10, "vote_average", "must be less or equal to 10")
	v.Check(movie.VoteCount >= 0, "votes_count", "must be a positive integer")

	v.Check(len(movie.Abstract) > 4096, "abstract", "must not be longer than 4096")

}
