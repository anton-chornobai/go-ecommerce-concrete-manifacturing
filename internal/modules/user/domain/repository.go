package domain

import (
	"context"
	"time"
)

type Repository interface {
	Signup(user *User) error
	SignupByEmail(ctx context.Context,user *User, verification_hash string, expires_at *time.Time) error
	LoginByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(number string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(id string) (*User, error)
	SaveVerificationCode(ctx context.Context, email, code string) error
	MarkUserVerified(ctx context.Context, email string) error
	IsAdmin(id string) (bool, error)
}
