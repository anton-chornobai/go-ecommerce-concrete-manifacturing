package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/contact/domain"
)

type UserContactRepository struct {
	DB *sql.DB
}

func (r *UserContactRepository) Save(ctx context.Context, userContact *domain.UserContact) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO user_contacts (id, name, email, number, message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		userContact.ID,
		userContact.Name,
		userContact.Email,
		userContact.Number,
		userContact.Message,
		userContact.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("insert user_contact: %w", err)
	}

	return nil
}

func (r *UserContactRepository) Delete(ctx context.Context, id string) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM user_contacts WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete user_contact: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user_contact, failed to check affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("no user_contact found with id %s", id)
	}

	return nil
}
