package database

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/base32"
	"encoding/json"
	"errors"
	"time"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
	ScopeChangeEmail    = "change-email"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Scope     string    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
	Data      TokenData `json:"-"`
}

type TokenData map[string]any

func (d TokenData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *TokenData) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &d)
}

func generateToken(userid int64, ttl time.Duration, scope string, tokenData TokenData) (*Token, error) {
	token := Token{
		UserId: userid,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
		Data:   tokenData,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return &token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m *TokenModel) New(userid int64, ttl time.Duration, scope string, tokenData TokenData) (*Token, error) {
	token, err := generateToken(userid, ttl, scope, tokenData)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m *TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope, data)
		VALUES($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope, token.Data}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&token.CreatedAt)
}

func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens 
		WHERE scope = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

func (m *TokenModel) ValidTokenAge(maxAge time.Duration, token *Token) bool {
	return time.Since(token.CreatedAt) < maxAge
}

func (m *TokenModel) GetByTokenHash(scope string, tokenHash []byte) (*Token, error) {
	query := `
		SELECT hash, user_id, scope, expiry, created_at, data FROM tokens
		WHERE scope = $1 AND hash = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	token := Token{}
	err := m.DB.QueryRowContext(ctx, query, scope, tokenHash).Scan(
		&token.Hash,
		&token.UserId,
		&token.Scope,
		&token.Expiry,
		&token.CreatedAt,
		&token.Data,
	)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
