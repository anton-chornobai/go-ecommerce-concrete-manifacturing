package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
	// "errors"
	// "github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersRepository struct {
	DB *sql.DB
}

func (o *OrdersRepository) Orders(ctx context.Context, limit int) ([]domain.Order, error) {
	rows, err := o.DB.QueryContext(ctx, `
		SELECT id, user_id, order_name, total, status, payment_status, discount, shipping_address, shipping_city, shipping_postal_code, created_at, updated_at FROM orders LIMIT=$1
	`, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []domain.Order

	for rows.Next() {
		var order domain.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderName,
			&order.Total,
			&order.Status,
			&order.PaymentStatus,
			&order.Discount,
			&order.ShippingAddress,
			&order.ShippingCity,
			&order.ShippingPostalCode,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
func (o *OrdersRepository) Create(ctx context.Context, order *domain.Order) error {
	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx: %w", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (
			user_id,
			order_name,
			total,
			status,
			payment_status,
			discount,
			shipping_address,
			shipping_city,
			shipping_postal_code,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id
	`,
		order.UserID,
		order.OrderName,
		order.Total,
		order.Status,
		order.PaymentStatus,
		order.Discount,
		order.ShippingAddress,
		order.ShippingCity,
		order.ShippingPostalCode,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_item (
				order_id,
				product_id,
				title,
				unit_price,
				type,
				quantity,
				color,
				material,
				height,
				width,
				thickness
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		`,
			order.ID,
			item.ProductID,
			item.Title,
			item.UnitPrice,
			item.Type,
			item.Quantity,
			item.Color,
			item.Material,
			item.Size.Height,
			item.Size.Width,
			item.Size.Thickness,
		)

		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (o *OrdersRepository) Delete(ctx context.Context, id int) error {
	res, err := o.DB.ExecContext(ctx, `DELETE FROM orders where id=$1`, id)

	if err != nil {
		return err
	}

	affectedRows, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("order with id %d not found", id)
	}

	return nil
}

// func (o *OrdersRepository) Update(ctx context.Context, id int) error {
// 	res, err := o.DB.ExecContext(ctx, `DELETE FROM users where id=$1`, id)

// }
