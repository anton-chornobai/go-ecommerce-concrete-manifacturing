package domain

import "context"

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	// Delete(ctx context.Context, id int) error
	// Edit(ctx context.Context, id int) error
	Get(ctx context.Context, limit int) ([]Order, error)
}
