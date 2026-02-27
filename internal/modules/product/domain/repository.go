package domain

import "context"

type Repository interface {
	Add(ctx context.Context, product *Product) (int, error)
	// Delete(ctx context.Context, id int) error 
	// Edit(ctx context.Context, id int)
	// Get()
}