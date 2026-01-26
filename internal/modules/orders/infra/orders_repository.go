package infra

import (
	"database/sql"
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
	// "errors"
	// "github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersRepository struct {
	DB *sql.DB
	Logger *slog.Logger
}

func (o *OrdersRepository) Orders(limit int) ([]domain.Order, error) {
	rows, err := o.DB.Query(`
	SELECT id, user_id, name, total_price, status, discount, created_at FROM orders LIMIT=?
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
			&order.Name,
			&order.Total,
			&order.Status,
			&order.Discount,
			&order.CreatedAt,
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

// func (o *OrdersRepository) OrderWithUserInfo(n int) (domain.Order, error) {
// 	var order domain.Order

// 	row := o.DB.QueryRow("SELECT id, user_id, user, product FROM orders WHERE user_id = ?", n)

// 	err := row.Scan(
// 		&order.ID,
// 		&order.UserId,
// 		&order.User,
// 		&order.Product,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return domain.Order{}, errors.New("no such row")
// 		}
// 		return domain.Order{}, err
// 	}

// 	return order, nil
// }

func (o *OrdersRepository) Save(order *domain.Order) error {
	result, err := o.DB.Exec(`
		INSERT INTO orders (user_id, name, total, status, discount, created_at) VALUES ( ?, ?, ?, ?, ?, ?);
	`, order.UserID, order.Name, order.Total, order.Status, order.Discount, order.CreatedAt)

	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := o.DB.Exec(`
			INSERT INTO order_items
			(order_id, product_id, title, unit_price, type, quantity, color, height, width, material, thickness)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			order.ID,
			item.ProductID,
			item.Title,
			item.UnitPrice,
			item.Type,
			item.Quantity,
			item.Color,
			item.Height,
			item.Width,
			item.Material,
			item.Thickness,
		)
		if err != nil {
			return err
		}
	}
	generatedId, err := result.LastInsertId()

	if err != nil {
		return err
	}

	order.ID = int(generatedId)

	return nil
}
