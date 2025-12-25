package domain

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID string `json:"id"`
	Number string `json:"number"`
	Role string `json:"role"`
	CreatedAt int `json:"created_at"`
	Email string `json:"email,omitempty"`
	Adress string `json:"adress,omitempty"`
	Name string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}	

type AuthenticationUserRequest struct {
	Number string `json:"number"`
}

type AuthenticationUserCreated struct {
	Number string `json:"number"`
	Role string `json:"role"`
	ID string `json:"id"`
}

type Claims struct {
	Number string `json:"number"`
	Role string `json:"role"`
	ID string `json:"id"`
	jwt.RegisteredClaims
}