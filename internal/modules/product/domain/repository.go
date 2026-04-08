package domain

import (
	"context"
)

type Repository interface {
	GetProducts(ctx context.Context, limit int, status *ProductStatus) ([]Product, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Add(ctx context.Context, product *Product) error
	DeleteByID(ctx context.Context, id int) error
	Update(ctx context.Context, id int, req ProductUpdate) error
}
