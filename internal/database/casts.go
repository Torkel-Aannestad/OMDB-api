package database

import (
	"context"
	"database/sql"
	"time"
)

type Cast struct {
	MovieID    int64     `json:"movie_id"`
	PersonID   int64     `json:"person_id"`
	JobID      int64     `json:"job_id"`
	Role       string    `json:"role"`
	Position   int32     `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type CastsModel struct {
	DB *sql.DB
}

func (m CastsModel) InsertCast(job *Job) error {
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
	RETURNING created_at, modified_at`

	return m.DB.QueryRowContext(ctx, query, job.Name).Scan(
		&job.ID,
		&job.CreatedAt,
		&job.ModifiedAt,
	)
}

func (m CastsModel) GetCastsByMovieID(movieID int64) ([]*Cast, error) {
	if movieID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
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
			&cast.MovieID,
			&cast.PersonID,
			&cast.JobID,
			&cast.Role,
			&cast.Position,
			&cast.CreatedAt,
			&cast.ModifiedAt,
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

func (m CastsModel) GetCastsByPersonID(personID int64) ([]*Cast, error) {
	if personID < 0 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT 
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
			&cast.MovieID,
			&cast.PersonID,
			&cast.JobID,
			&cast.Role,
			&cast.Position,
			&cast.CreatedAt,
			&cast.ModifiedAt,
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

func (m CastsModel) DeleteCast(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM casts WHERE movie_id = $1 AND person_id = $2 AND job_id = $3
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
