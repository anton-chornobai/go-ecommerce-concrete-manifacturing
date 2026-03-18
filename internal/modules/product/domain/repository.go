package domain

import (
	"context"
)

type Repository interface {
	GetWithLimit(ctx context.Context, limit int) ([]Product, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Add(ctx context.Context, product *Product) error
	RemoveByID(ctx context.Context, id int) error
	Update(ctx context.Context, id int, req ProductUpdate) error
	
}
