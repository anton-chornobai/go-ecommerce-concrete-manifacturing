package domain

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const phoneNumberLength = 9

type User struct {
	ID                    string
	Role                  string
	Address               string
	Name                  string
	Surname               string
	Email                 string
	Password              string
	VerificationHash      string
	IsVerified            bool
	Number                *string
	CreatedAt             time.Time
	VerificationExpiresAt *time.Time
}
type Claims struct {
	Number string `json:"number"`
	Role   string `json:"role"`
	ID     string `json:"id"`
	jwt.RegisteredClaims
}

func CreateUser(number string) (*User, error) {
	if number == "" || len(number) < phoneNumberLength {
		return nil, errors.New("invalid phone number")
	}

	return &User{
		ID:     uuid.NewString(),
		Role:   "customer",
		Number: &number,
	}, nil
}

func CreateUserWithEmail(email, password string) *User {
	return &User{
		ID:       uuid.NewString(),
		Role:     "customer",
		Email:    email,
		Password: password,
	}
}
