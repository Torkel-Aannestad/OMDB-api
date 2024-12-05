package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type PeopleLink struct {
	ID         int64     `json:"id"`
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
	RETURNING id, created_at, modified_at`

	args := []any{
		peopleLink.Source,
		peopleLink.Key,
		peopleLink.PersonID,
		peopleLink.Language,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&peopleLink.ID,
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
		id,
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
			&personLink.ID,
			&personLink.Key,
			&personLink.Source,
			&personLink.PersonID,
			&personLink.Language,
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

func (m PeopleLinkModel) Delete(id int64) error {

	stmt := `
		DELETE FROM people_links WHERE id = $1;
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

func ValidatePeopleLink(v *validator.Validator, PeopleLink *PeopleLink) {
	v.Check(PeopleLink.Key != "", "key", "must be provided")
	v.Check(validator.PermittedValue(PeopleLink.Key, "wikidata", "wikipedia", "imdbperson"), "source", "must not be of the following values 'wikidata', 'wikipedia' or 'imdbperson'")

	v.Check(PeopleLink.Source != "", "source", "must be provided")
	v.Check(len(PeopleLink.Source) <= 500, "source", "must not be more than 500 bytes long")
}
