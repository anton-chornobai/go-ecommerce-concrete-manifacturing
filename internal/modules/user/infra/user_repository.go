package infra

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *UserRepository) SignupByEmail(ctx context.Context, user *domain.User) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO users (id, role, email, password) VALUES ($1, $2, $3, $4)`, user.ID, user.Role, user.Email, user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) LoginByEmail(ctx context.Context, email, password string) (*domain.User, error ) {
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

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	row := r.DB.QueryRow(
		"SELECT id, number, name, surname, role, email, created_at, address FROM users WHERE email=?", email,
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
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
