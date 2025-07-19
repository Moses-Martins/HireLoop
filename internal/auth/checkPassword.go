package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordByte := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil

}

func CheckPasswordHash(password, hash string) error {
	hashedByte := []byte(hash) 

	err := bcrypt.CompareHashAndPassword(hashedByte, []byte(password))
	if err != nil {
		return err
	}

	return nil
}