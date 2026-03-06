package application

import (
	"context"
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductService struct {
	repo domain.Repository
	log  *slog.Logger
}

func NewProductService(repo domain.Repository) (*ProductService, error) {
	return &ProductService{repo: repo}, nil
}

func (p *ProductService) Add(ctx context.Context, input domain.Product) error {
	product, err := domain.NewProduct(
		input.Price,
		input.Title,
		input.Type,
		input.ImageURL,
		input.Color,
		input.StockQuantity,
		input.Description,
		input.Weight,
		input.Rating,
		input.Size,
	)

	if err != nil {
		return err
	}

	err = p.repo.Add(ctx, product)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductService) DeleteByID(ctx context.Context, id int) error {
	err := p.repo.RemoveByID(ctx, id)

	if err != nil {
		return err;
	}

	return nil
}

func (p *ProductService) Edit(ctx context.Context) error {
	return nil
}
