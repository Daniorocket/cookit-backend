package lib

import (
	"golang.org/x/crypto/bcrypt"
)

var cost = 8

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func CompareHashAndPassword(passwordFromDB, passwordFromJSON []byte) error {
	if err := bcrypt.CompareHashAndPassword(passwordFromDB, passwordFromJSON); err != nil {
		return err
	}
	return nil
}
