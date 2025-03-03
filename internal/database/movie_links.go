package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/OMDB-api/internal/validator"
)

type MovieLink struct {
	ID         int64     `json:"id"`
	Source     string    `json:"source"`
	Key        string    `json:"key"`
	MovieID    int64     `json:"movie_id"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type MovieLinkModel struct {
	DB *sql.DB
}

func (m MovieLinkModel) Insert(movieLink *MovieLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
	INSERT INTO movie_links (
		source,  
		key,
		movie_id,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, modified_at`

	args := []any{
		movieLink.Source,
		movieLink.Key,
		movieLink.MovieID,
		movieLink.Language,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&movieLink.ID,
		&movieLink.CreatedAt,
		&movieLink.ModifiedAt,
	)
}

func (m MovieLinkModel) Get(movieID int64) ([]*MovieLink, error) {
	if movieID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		id,
		source,  
		key,
		movie_id,
		language
	FROM movie_links
	WHERE movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieLinks := []*MovieLink{}

	for rows.Next() {
		var personLink MovieLink

		err := rows.Scan(
			&personLink.ID,
			&personLink.Key,
			&personLink.Source,
			&personLink.MovieID,
			&personLink.Language,
		)
		if err != nil {
			return nil, err
		}

		movieLinks = append(movieLinks, &personLink)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return movieLinks, nil
}

func (m MovieLinkModel) Delete(Id int64) error {
	stmt := `
		DELETE FROM movie_links WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, Id)
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

func ValidateMovieLink(v *validator.Validator, MovieLink *MovieLink) {
	v.Check(MovieLink.Source != "", "source", "must be provided")
	v.Check(len(MovieLink.Source) <= 500, "source", "must not be more than 500 bytes long")

	v.Check(MovieLink.Key != "", "key", "must be provided")
	v.Check(validator.PermittedValue(MovieLink.Key, "wikidata", "wikipedia", "imdbperson"), "source", "must not be of the following values 'wikidata', 'wikipedia' or 'imdbperson'")

}
