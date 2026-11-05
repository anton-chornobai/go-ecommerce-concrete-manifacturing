package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetProducts(ctx context.Context, limit int, status *ProductStatus) ([]Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	Create(ctx context.Context, product *Product) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Patch(ctx context.Context, id uuid.UUID, req *ProductPatch) error
}
