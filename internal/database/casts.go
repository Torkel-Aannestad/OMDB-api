package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Cast struct {
	ID         int64     `json:"id"`
	MovieID    int64     `json:"movie_id"`
	PersonID   int64     `json:"person_id"`
	JobID      int64     `json:"job_id"`
	Role       string    `json:"role"`
	Position   int32     `json:"position"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

type CastsModel struct {
	DB *sql.DB
}

func (m CastsModel) Insert(cast *Cast) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO casts (
		movie_id,
		person_id,
		job_id,
		role,
		position
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, modified_at`

	args := []any{cast.MovieID, cast.PersonID, cast.JobID, cast.Role, cast.Position}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&cast.ID,
		&cast.CreatedAt,
		&cast.ModifiedAt,
	)
}

func (m CastsModel) GetByMovieID(movieID int64) ([]*Cast, error) {
	if movieID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
		id,
		movie_id,
		person_id,
		job_id,
		role,
		position
	FROM casts
	WHERE movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	casts := []*Cast{}

	for rows.Next() {
		var cast Cast

		err := rows.Scan(
			&cast.ID,
			&cast.MovieID,
			&cast.PersonID,
			&cast.JobID,
			&cast.Role,
			&cast.Position,
		)
		if err != nil {
			return nil, err
		}
		casts = append(casts, &cast)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return casts, nil
}

func (m CastsModel) GetByPersonID(personID int64) ([]*Cast, error) {
	if personID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
		id,
		movie_id,
		person_id,
		job_id,
		role,
		position
	FROM casts
	WHERE person_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	casts := []*Cast{}

	for rows.Next() {
		var cast Cast

		err := rows.Scan(
			&cast.ID,
			&cast.MovieID,
			&cast.PersonID,
			&cast.JobID,
			&cast.Role,
			&cast.Position,
		)
		if err != nil {
			return nil, err
		}
		casts = append(casts, &cast)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return casts, nil
}

func (m CastsModel) Delete(id int64) error {

	stmt := `
		DELETE FROM casts WHERE id = $1;
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

func ValidateCast(v *validator.Validator, cast *Cast) {
	v.Check(cast.MovieID != 0, "movie_id", "must be provided")
	v.Check(cast.MovieID > 0, "movie_id", "must be a positive number")

	v.Check(cast.PersonID != 0, "person_id", "must be provided")
	v.Check(cast.PersonID > 0, "person_id", "must be a positive number")

	v.Check(cast.JobID != 0, "job_id", "must be provided")
	v.Check(cast.JobID > 0, "job_id", "must be a positive number")

	v.Check(len(cast.Role) <= 250, "role", "must not be longer than 250 characters")

	v.Check(cast.Position != 0, "position", "must be provided")
	v.Check(cast.Position > 0, "position", "must be a positive number")

}
