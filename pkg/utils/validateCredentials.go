package utils

import (
	"errors"
	"regexp"
)

func ValidatePasswordAndEmail(email, password string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}