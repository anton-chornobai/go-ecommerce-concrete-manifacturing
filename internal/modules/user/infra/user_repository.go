package infra

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Signup(user *domain.User) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, number, role) VALUES (?, ?, ?)`,
		user.ID,
		user.Number,
		user.Role,
	)

	if err != nil {

	}

	return err
}

func (r *UserRepository) SignupByEmail(ctx context.Context, user *domain.User, verificationHash string, expiresAt *time.Time) error {
	_, err := r.DB.ExecContext(
		ctx,
		`INSERT INTO users (id, role, email, password, verification_hash, verification_expires_at) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID,
		user.Role,
		user.Email,
		user.Password,
		verificationHash,
		expiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) LoginByEmail(ctx context.Context, email, password string) (*domain.User, error) {
	var user domain.User

	row := r.DB.QueryRowContext(ctx, `SELECT password, role, name FROM users WHERE email=$1`, email)
	err := row.Scan(
		&user.Password,
		&user.Role,
		&user.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no such user with email %s", email)
		}
		return nil, fmt.Errorf("error scanning password: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) SaveVerificationCode(ctx context.Context, email, hashed_code string) error {
	res, err := r.DB.ExecContext(ctx,
		`UPDATE users SET verification_code_hash=$1, verification_expires_at = NOW() + INTERVAL '10 minutes' WHERE email=$2;`,
		hashed_code, email)

	if err != nil {
		return fmt.Errorf("failed to update verification code: %w", err)
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if affectedRow == 0 {
		return fmt.Errorf("no user found with email %s", email)
	}

	return nil
}

func (r *UserRepository) GetByPhone(number string) (*domain.User, error) {
	var user domain.User
	row := r.DB.QueryRow(
		"SELECT id, number, name, surname, role, email, created_at, address FROM users WHERE number=?", number)

	err := row.Scan(
		&user.ID,
		&user.Number,
		&user.Name,
		&user.Surname,
		&user.Role,
		&user.Email,
		&user.CreatedAt,
		&user.Address,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	row := r.DB.QueryRowContext(
		ctx,
		`SELECT id, number, name, surname, role, email, created_at, address, verified, verification_hash, verification_expires_at
		 FROM users
		 WHERE email = $1`,
		email,
	)

	err := row.Scan(
		&user.ID,
		&user.Number,
		&user.Name,
		&user.Surname,
		&user.Role,
		&user.Email,
		&user.CreatedAt,
		&user.Address,
		&user.IsVerified,
		&user.VerificationHash,
		&user.VerificationExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) MarkUserVerified(ctx context.Context, email string) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET verified = TRUE,
		    verification_hash = '',
		    verification_expires_at = NULL
		WHERE email = $1
		  AND verified = FALSE
	`, email)
	if err != nil {
		return fmt.Errorf("failed to mark user verified: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no unverified user found with email %s", email)
	}

	return nil
}
