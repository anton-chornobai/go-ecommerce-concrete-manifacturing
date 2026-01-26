package domain

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User struct {
	ID        string `json:"id"`
	Number    string `json:"number"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	Email     *string `json:"email,omitempty"`
	Password  string `json:"password"`
	Address   string `json:"adress,omitempty"`
	Name      string `json:"name,omitempty"`
	Surname   string `json:"surname,omitempty"`
}

type AuthenticationUserRequest struct {
	Number string `json:"number"`
}

type AuthenticationUserCreated struct {
	Number string `json:"number"`
	Role   string `json:"role"`
	ID     string `json:"id"`
}

type Claims struct {
	Number string `json:"number"`
	Role   string `json:"role"`
	ID     string `json:"id"`
	jwt.RegisteredClaims
}

type RegisterResult struct {
	User  UserCreated
	Token string
}

const phoneNumberLength = 9

type UserCreated struct {
	ID     string
	Role   string
	Number string
}

type UserCreatedWithEmail struct {
	ID     string
	Role   string
	Email string
	Password string
}

func CreateUser(number string) (*UserCreated, error) {
	if number == "" || len(number) < phoneNumberLength {
		return nil, errors.New("invalid phone number")
	}

	return &UserCreated{
		ID:     uuid.NewString(),
		Role:   "user",
		Number: number,
	}, nil
}

func CreateUserWithEmail(email, password string) (*UserCreatedWithEmail) {
	return &UserCreatedWithEmail{
		ID:     uuid.NewString(),
		Role:   "user",
		Email: email,
		Password: password,
	}
}