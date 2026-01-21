package infra

import (
	"database/sql"
	"log"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user domain.UserCreated) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, number, role) VALUES (?, ?, ?)`,
		user.ID,
		user.Number,
		user.Role,
	)

	if err != nil {
		log.Printf("RegisterUser error: %+v", err)
	}

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
		return domain.User{}, err
	}

	return user, nil
}