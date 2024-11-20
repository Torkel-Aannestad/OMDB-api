package auth

import (
	"errors"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePlaintextPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "can't be blank")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
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
