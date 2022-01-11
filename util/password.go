package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("%s: %w", "Couldn't hash password", err)
	}
	return string(hashpassword), nil
}

func CheckPassword(password, hashpassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
	return err
}
