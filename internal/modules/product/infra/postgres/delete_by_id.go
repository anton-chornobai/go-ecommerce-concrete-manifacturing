package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (p *ProductRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE order_item SET product_id = NULL WHERE product_id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to nullify product references: %w", err)
	}

	res, err := tx.ExecContext(ctx,
		`DELETE FROM products WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
