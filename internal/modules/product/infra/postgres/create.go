package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/jackc/pgx/v5/pgconn"
)


func (p *ProductRepository) Create(ctx context.Context, product *domain.Product) error {

	var sizeWidth *int
	var sizeHeight *int

	if product.Size != nil {
		sizeWidth = &product.Size.Width
		sizeHeight = &product.Size.Height
	}

	var productID int
	err := p.DB.QueryRowContext(ctx, `
		INSERT INTO products (
			price, title, type, status, color,
			description, stock_quantity, weight_grams, rating,
			size_width, size_height
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11
		) RETURNING id`,
		product.Price,
		product.Title,
		product.Type,
		product.Status,
		product.Color,
		product.Description,
		product.StockQuantity,
		product.Weight,
		product.Rating,
		sizeWidth,
		sizeHeight,
	).Scan(&productID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrTitleAlreadyExists
		}
		return fmt.Errorf("не вдалося додати продукт: %w", err)
	}

	for i, image := range product.ImageURLs {
		_, err := p.DB.ExecContext(ctx, `
			INSERT INTO product_image (id, product_id, url, position)
			VALUES ($1, $2, $3, $4)`,
			image.ID,
			productID,
			image.Key,
			i,
		)
		if err != nil {
			return fmt.Errorf("не вдалося зберегти зображення %s: %w", image.Key, err)
		}
	}

	return nil
}
