package domain

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("користувача не знайдено")
	ErrUnauthorized = errors.New("недостатньо прав")
)

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

func CreateUserWithEmail(email, password string) *User {
	return &User{
		ID:       uuid.NewString(),
		Role:     "customer",
		Email:    email,
		Password: password,
	}
}
