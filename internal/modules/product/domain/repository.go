package domain

import "context"

type Repository interface {
	Add(ctx context.Context, product *Product)  error
	RemoveByID(ctx context.Context, id int) error 
	// Edit(ctx context.Context, id int)
	// Get()
}