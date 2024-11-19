package auth

import (
	"errors"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func GenerateHashFromPlaintext(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func PasswordMatches(plaintextPassword string, hashedPassord []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassord, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
