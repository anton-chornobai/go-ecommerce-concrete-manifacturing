package infra

import (
	"database/sql"
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UserRepository struct {
	DB *sql.DB
	Logger *slog.Logger
}

func (r *UserRepository) Create(user domain.UserCreated) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, number, role) VALUES (?, ?, ?)`,
		user.ID,
		user.Number,
		user.Role,
	)

	if err != nil {
		r.Logger.Error("Failed to insert user", "err", err)
	}

	r.Logger.Info(
		"User Created",
		slog.String("user_id", user.ID),
		slog.String("user_role", user.Role),
	)
	return err 
}

func (r *UserRepository) GetByPhone(number string) (domain.User, error) {
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
		r.Logger.Warn("Failed to scan user row by phone", "err:", err, "phone", number)
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (domain.User, error) {
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

		
		r.Logger.Warn("Failed to scan user row by email", "err:", err)
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) SignUpByEmail(user *domain.UserCreatedWithEmail) error  {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, role, email, password) VALUES (?, ?, ?, ?)`, user.ID, user.Role, user.Email, user.Password,
	)

	if err != nil {
		r.Logger.Warn("Failed to scan user row by email", "err:", err)
		return  err;
	}

	return  nil
}
