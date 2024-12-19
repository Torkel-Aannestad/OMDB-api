package auth

import (
	"crypto/rand"
	"errors"

	"github.com/Torkel-Aannestad/MovieMaze/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePlaintextPassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "can't be blank")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
	v.Check(validator.NotIn(password, CommonPasswords...), "password", "password is commonly used and will not be accepted")
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

type ParamsArgon2 struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParamsArgon2 = &ParamsArgon2{
	Memory:      64 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
