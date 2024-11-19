package database

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	Activated    bool      `json:"activated"`
	Version      int       `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	args := []any{user.Name, user.Email, user.PasswordHash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		} else {
			return err
		}
	}
	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	user := User{}
	query := `
		SELECT 
		id, created_at, name, email, password_hash, activated, version 
		FROM users
		WHERE email = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, ErrRecordNotFound
		} else {
			return &user, err
		}
	}
	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET name = $2, email = $3, password_hash = $4, activated = $5, version = version + 1
        WHERE id = $1 AND version = $6
        RETURNING version`

	args := []any{
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Activated,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	hashArray := sha256.Sum256([]byte(tokenPlaintext))
	hash := hashArray[:]

	query := `
		SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3
	`
	args := []any{hash, tokenScope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := User{}

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email can't be empty")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "email must be valid email")
}

func ValidatePlaintextPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "can't be blank")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "can't be blank")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.PasswordHash == nil {
		panic("missing password hash for user")
	}
}
