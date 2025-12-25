package infra

import (
	"database/sql"

	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) FindByPhone(number int) (domain.User, error) {
	var user domain.User 
	row := r.DB.QueryRow("SELECT id, number, name, surname, role, emain, created_at, address FROM users WHERE number=?", number)
	err := row.Scan(&user.ID, &user.Name, &user.Role)

	if err != nil {
		return domain.User{}, err
	}
	return user, err
}


func (r *UserRepository) Create(user domain.User) error {
	_, err := r.DB.Exec(
		`INSERT INTO users (id, number, role) VALUES (?, ?, ?)`,
		user.ID,
		user.Number,
		user.Role,
	)

	return err
}