package infra

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {}

func (h *PasswordHasher) HashPassword(password string) ([]byte, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {		
		return nil, errors.New("123")
	}

	return hashedPassword, nil
}

func (h *PasswordHasher) CompareHashAndPassword(hashedPassword, password  string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return errors.New("passwords hash do not match")
	}

	return nil
}