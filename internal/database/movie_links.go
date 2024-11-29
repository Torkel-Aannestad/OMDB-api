package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type MovieLink struct {
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO movie_links (
		source,  
		key,
		movie_id,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at, modified_at`

	args := []any{
		movieLink.Source,
		movieLink.Key,
		movieLink.MovieID,
		movieLink.Language,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&movieLink.CreatedAt,
		&movieLink.ModifiedAt,
	)
}

func (m MovieLinkModel) GetMovieLinks(personID int64) ([]*MovieLink, error) {
	if personID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
	source,  
	key,
	movie_id,
	language
	FROM movie_links
	WHERE movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieLinks := []*MovieLink{}

	for rows.Next() {
		var personLink MovieLink

		err := rows.Scan(
			&personLink.Key,
			&personLink.Source,
			&personLink.MovieID,
			&personLink.Language,
			&personLink.CreatedAt,
			&personLink.ModifiedAt,
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

func (m MovieLinkModel) Delete(id int64) error {
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

func ValidateMovieLink(v *validator.Validator, MovieLink *MovieLink) {
	v.Check(MovieLink.Key != "", "key", "must be provided")
	v.Check(len(MovieLink.Key) <= 500, "key", "must not be more than 500 bytes long")

	v.Check(MovieLink.Source != "", "source", "must be provided")
	v.Check(len(MovieLink.Source) <= 500, "source", "must not be more than 500 bytes long")
	v.Check(validator.PermittedValue(MovieLink.Source, "wikidata", "wikipedia", "imdbperson"), "source", "must not be of the following values 'wikidata', 'wikipedia' or 'imdbperson'")
}
