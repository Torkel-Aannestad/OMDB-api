package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

type Job struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

func (m CastsModel) InsertJob(job *Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO jobs (
		name
	)
	VALUES ($1)
	RETURNING id, created_at, modified_at`

	return m.DB.QueryRowContext(ctx, query, job.Name).Scan(
		&job.ID,
		&job.CreatedAt,
		&job.ModifiedAt,
	)
}

func (m CastsModel) GetJobById(id int64) (*Job, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	job := Job{}

	query := `
		SELECT 
			id,  
			name,
			created_at,
			modified_at
		FROM jobs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Name,
		&job.CreatedAt,
		&job.ModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}

	return &job, nil
}

func (m CastsModel) DeleteJob(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM jobs WHERE id = $1
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

func ValidateJob(v *validator.Validator, job *Category) {
	v.Check(job.Name != "", "name", "must be provided")
	v.Check(len(job.Name) <= 500, "name", "must not be more than 500 bytes long")
}
