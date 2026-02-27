package application

import (
	"context"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductService struct {
	repo domain.Repository
}

func NewProductService(repo domain.Repository) (*ProductService, error) {
	return &ProductService{repo: repo}, nil
}

func (p *ProductService) Add(ctx context.Context, input domain.Product) (int, error) {
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
		return -1, err
	}

	createdID, err := p.repo.Add(ctx, product)

	if err != nil {
		return -1, err
	}

	return createdID, nil
}

func (p *ProductService) Remove(ctx context.Context) error {
	return nil

}

func (p *ProductService) Edit(ctx context.Context) error {
	return nil
}
