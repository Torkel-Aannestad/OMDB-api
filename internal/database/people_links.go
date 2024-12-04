package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type PeopleLink struct {
	Source     string    `json:"source"`
	Key        string    `json:"key"`
	PersonID   int64     `json:"person_id"`
	Language   string    `json:"language"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type PeopleLinkModel struct {
	DB *sql.DB
}

func (m PeopleLinkModel) Insert(peopleLink *PeopleLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO people_links (
		source,  
		key,
		person_id,
		language
	)
	VALUES ($1, $2, $3, $4)
	RETURNING created_at, modified_at`

	args := []any{
		peopleLink.Source,
		peopleLink.Key,
		peopleLink.PersonID,
		peopleLink.Language,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&peopleLink.CreatedAt,
		&peopleLink.ModifiedAt,
	)
}

func (m PeopleLinkModel) Get(personID int64) ([]*PeopleLink, error) {
	if personID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
	source,  
	key,
	person_id,
	language
	FROM people_links
	WHERE person_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	peopleLinks := []*PeopleLink{}

	for rows.Next() {
		var personLink PeopleLink

		err := rows.Scan(
			&personLink.Key,
			&personLink.Source,
			&personLink.PersonID,
			&personLink.Language,
			&personLink.CreatedAt,
			&personLink.ModifiedAt,
		)
		if err != nil {
			return nil, err
		}

		peopleLinks = append(peopleLinks, &personLink)

	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return peopleLinks, nil
}

func (m PeopleLinkModel) Delete(personID int64, language, key string) error {

	stmt := `
		DELETE FROM people_links WHERE person_id = $1 AND language = $2 AND key = $3;
	`
	args := []any{personID, language, key}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, stmt, args...)
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

func ValidatePeopleLink(v *validator.Validator, PeopleLink *PeopleLink) {
	v.Check(PeopleLink.Key != "", "key", "must be provided")
	v.Check(len(PeopleLink.Key) <= 500, "key", "must not be more than 500 bytes long")

	v.Check(PeopleLink.Source != "", "source", "must be provided")
	v.Check(len(PeopleLink.Source) <= 500, "source", "must not be more than 500 bytes long")
	v.Check(validator.PermittedValue(PeopleLink.Source, "wikidata", "wikipedia", "imdbperson"), "source", "must not be of the following values 'wikidata', 'wikipedia' or 'imdbperson'")
}
