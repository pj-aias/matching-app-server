package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) ([]byte, error) {
	// bcrypt won't work correctly if the password length is > 72
	if len(password) > 72 {
		err := errors.New("Password length must be less than 72 bytes.")
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
