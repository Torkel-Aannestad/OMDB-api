package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/OMDB-api/internal/validator"
)

type Job struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
	Version    int32     `json:"version"`
}

type JobsModel struct {
	DB *sql.DB
}

func (m JobsModel) Insert(job *Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
	INSERT INTO jobs (
		name
	)
	VALUES ($1)
	RETURNING id, created_at, modified_at, version`

	return m.DB.QueryRowContext(ctx, query, job.Name).Scan(
		&job.ID,
		&job.CreatedAt,
		&job.ModifiedAt,
		&job.Version,
	)
}

func (m JobsModel) Get(id int64) (*Job, error) {
	if id < 0 {
		return nil, ErrRecordNotFound
	}
	job := Job{}

	query := `
		SELECT 
			id,  
			name,
			created_at,
			modified_at,
			version
		FROM jobs
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Name,
		&job.CreatedAt,
		&job.ModifiedAt,
		&job.Version,
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

func (m JobsModel) Update(job *Job) error {
	query := `
	UPDATE jobs
	SET 
		name = $3,
		modified_at = NOW(),
		version = version + 1
	WHERE id = $1 and version = $2
	RETURNING version`

	args := []any{
		&job.ID,
		&job.Version,
		&job.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&job.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		} else {
			return err
		}
	}
	return nil
}

func (m JobsModel) Delete(id int64) error {
	if id < 0 {
		return ErrRecordNotFound
	}

	stmt := `
		DELETE FROM jobs WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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

func ValidateJob(v *validator.Validator, job *Job) {
	v.Check(job.Name != "", "name", "must be provided")
	v.Check(len(job.Name) <= 500, "name", "must not be more than 500 bytes long")
}
