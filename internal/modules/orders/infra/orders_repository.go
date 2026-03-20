package infra

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersRepository struct {
	DB *sql.DB
}
func (o *OrdersRepository) Get(ctx context.Context, limit int) ([]domain.Order, error) {

	rows, err := o.DB.QueryContext(ctx, `
	SELECT 
	    o.id,
		o.customer_name,
		o.customer_number,
	    o.user_id,
	    o.order_name,
	    o.total,
	    o.status,
	    o.payment_status,
	    o.discount,
	    o.shipping_address,
	    o.shipping_city,
	    o.shipping_postal_code,
	    o.created_at,
	    o.updated_at,

	    oi.id,
	    oi.title,
	    oi.type,
	    oi.color,
	    oi.material,
	    oi.order_id,
	    oi.product_id,
	    oi.quantity,
	    oi.unit_price,
	    oi.height,
	    oi.width,
	    oi.thickness

	FROM orders o
	LEFT JOIN order_item oi 
	ON o.id = oi.order_id
	ORDER BY o.id
	LIMIT $1
	`, limit)

	if err != nil {
		return nil, err
	}

	orderMap := make(map[int]*domain.Order)

	for rows.Next() {

		var order domain.Order
		var item domain.OrderItem
		var size domain.Size

		err := rows.Scan(
			&order.ID,
			&order.CustomerName,
			&order.CustomerNumber,
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

			&item.ID,
			&item.Title,
			&item.Type,
			&item.Color,
			&item.Material,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&size.Height,
			&size.Width,
			&size.Thickness,
		)

		if err != nil {
			return nil, err
		}

		item.Size = size

		existingOrder, ok := orderMap[order.ID]

		if !ok {
			order.Items = []domain.OrderItem{}
			orderMap[order.ID] = &order
			existingOrder = &order
		}

		if item.ID.String() != "" {
			existingOrder.Items = append(existingOrder.Items, item)
		}
	}

	var orders []domain.Order

	for _, o := range orderMap {
		orders = append(orders, *o)
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
			customer_name,
			customer_number,
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
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id
	`,
		order.UserID,
		order.CustomerName,
		order.CustomerNumber,
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
