package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Trailer struct {
	TrailerID  int64     `json:"trailer_id"`
	Key        string    `json:"key"`
	MovieID    int64     `json:"movie_id"`
	Language   string    `json:"language"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type TrailersModel struct {
	DB *sql.DB
}

func (m TrailersModel) Insert(trailer *Trailer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO trailers (
		movie_id,
		source,  
		key,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING trailer_id, created_at, modified_at`

	args := []any{
		trailer.MovieID,
		trailer.Source,
		trailer.Key,
		trailer.Language,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&trailer.TrailerID,
		&trailer.CreatedAt,
		&trailer.ModifiedAt,
	)
}

func (m TrailersModel) GetTrailers(MovieID int64) ([]*Trailer, error) {
	if MovieID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
			trailer_id,
			source,  
			key,
			movie_id,
			language
		FROM trailers
		WHERE movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, MovieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trailers := []*Trailer{}

	for rows.Next() {
		var trailer Trailer

		err := rows.Scan(
			&trailer.MovieID,
			&trailer.Key,
			&trailer.Source,
			&trailer.Language,
			&trailer.CreatedAt,
			&trailer.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		trailers = append(trailers, &trailer)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return trailers, nil
}

func (m TrailersModel) Delete(trailerID int64) error {
	if trailerID < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM trailers WHERE trailer_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, trailerID)
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

func ValidateTrailer(v *validator.Validator, Trailer *Trailer) {
	v.Check(Trailer.Key != "", "key", "must be provided")
	v.Check(len(Trailer.Key) <= 500, "key", "must not be more than 500 bytes long")

	v.Check(Trailer.Source != "", "source", "must be provided")
	v.Check(len(Trailer.Source) <= 500, "source", "must not be more than 500 bytes long")
	v.Check(validator.PermittedValue(Trailer.Source, "youtube", "vimeo"), "source", "must not be of the following values 'youtube' or 'vimeo'")
}
