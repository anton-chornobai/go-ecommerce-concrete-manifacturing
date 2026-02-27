package infra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductRepository struct {
	DB *sql.DB
}

func (p *ProductRepository) Add(ctx context.Context, product *domain.Product) (int, error) {
	var id int

	var sizeWidth *int
	var sizeHeight *int

	// Size could be nil so accessing the fields that doesnt exist would panic,
	// therefore we check

	if product.Size != nil {
		sizeWidth = &product.Size.Width
		sizeHeight = &product.Size.Height
	}
	err := p.DB.QueryRowContext(ctx, `
		INSERT INTO product (
			price,
			title,
			product_type,
			image_url,
			color,
			description,
			stock_quantity,
			weight_grams,
			rating,
			size_width,
			size_height
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		RETURNING id;
	`, product.Price,
		product.Title,
		product.Type,
		product.ImageURL,
		product.Color,
		product.Description,
		product.StockQuantity,
		product.Weight,
		product.Rating,
		sizeWidth,
		sizeHeight,
	).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("couldnt add new product: %w", err)
	}

	return id, nil
}

func (p *ProductRepository) Remove(ctx context.Context) error {

	return nil
}

func (p *ProductRepository) Edit(ctx context.Context) error {

	return nil
}
