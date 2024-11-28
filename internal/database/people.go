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
	INSERT INTO people (
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

func (m PeopleModel) Update(person *Person) error {
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
		&person.ID,
		&person.Version,
		&person.Name,
		&person.Birthday,
		&person.Deathday,
		&person.Gender,
		&person.Aliases,
		&person.ModifiedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&person.Version)
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

func ValidatePeople(v *validator.Validator, person *Person) {
	v.Check(person.Name != "", "name", "must be provided")
	v.Check(len(person.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(person.Birthday.IsZero(), "date", "must be provided")
	v.Check(person.Birthday.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(person.Birthday.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(person.Deathday.Year() >= 1888, "date", "must be greater than year 1888")
	v.Check(person.Deathday.Compare(time.Now()) < 1, "date", "must not be in the future")

	v.Check(person.Gender != "", "gender", "must be provided")
	v.Check(person.Gender != "male" && person.Gender != "female" && person.Gender != "non-binary" && person.Gender != "unknown", "Gender", "must be one of the following values: male, female, non-binary, unknown")

	v.Check(validator.Unique(person.Aliases), "aliases", "must not contain duplicate values")

}
