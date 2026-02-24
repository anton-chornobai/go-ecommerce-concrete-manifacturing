package domain

import "context"

type Repository interface {
	Signup(user *User) error
	SignupByEmail(ctx context.Context,user *User) error
	LoginByEmail(ctx context.Context, email, password string) (*User, error)
	GetByPhone(number string) (*User, error)
	GetByEmail(email string) (*User, error)
}
