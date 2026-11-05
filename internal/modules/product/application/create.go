package application

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type CreateProductInput struct {
	Price         int
	Title         string
	Type          string
	Status        string
	Color         *string
	Description   *string
	StockQuantity *int
	Weight        *int
	Rating        *int
	SizeWidth     *int
	SizeHeight    *int
}

type FileInput struct {
	Reader      io.Reader
	ContentType string
	Filename    string
}

func (p *ProductService) Create(
	ctx context.Context,
	input CreateProductInput,
	fileInput []FileInput,
) (domain.Product, error) {
	var size *domain.Size

	if input.SizeWidth != nil && input.SizeHeight != nil {
		size = &domain.Size{
			Width:  *input.SizeWidth,
			Height: *input.SizeHeight,
		}
	}

	newProduct, err := domain.NewProduct(
		input.Price,
		input.Title,
		input.Type,
		input.Color,
		domain.ProductStatus(input.Status),
		nil,
		input.StockQuantity,
		input.Description,
		input.Weight,
		input.Rating,
		size,
	)
	if err != nil {
		return domain.Product{}, err
	}

	images, err := p.uploadImages(ctx, fileInput)
	if err != nil {
		return domain.Product{}, err
	}

	newProduct.ImageURLs = images

	if err := p.repo.Create(ctx, newProduct); err != nil {
		keys := make([]string, 0, len(images))
		for _, image := range images {
			keys = append(keys, image.Key)
		}

		cleanupErr := p.cleanupImageKeys(keys)
		return domain.Product{}, errors.Join(
			fmt.Errorf("create product: %w", err),
			cleanupErr,
		)
	}

	return *newProduct, nil
}
