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

func (p *ProductRepository) Add(ctx context.Context, product *domain.Product) error {
	var id int

	var sizeWidth *int
	var sizeHeight *int

	if product.Size != nil {
		sizeWidth = &product.Size.Width
		sizeHeight = &product.Size.Height
	}
	err := p.DB.QueryRowContext(ctx, `
		INSERT INTO products (
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
		return fmt.Errorf("couldnt add new product: %w", err)
	}

	return nil
}

func (p *ProductRepository) RemoveByID(ctx context.Context, id int) error {
	res, err := p.DB.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("couldnt exec deletion: %w", err)
	}

	affectedRow, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("couldnt check if rows were affected: %w", err)
	}

	if affectedRow == 0 {
		return fmt.Errorf("row not found")
	}

	return nil
}

func (p *ProductRepository) Edit(ctx context.Context) error {

	return nil
}
