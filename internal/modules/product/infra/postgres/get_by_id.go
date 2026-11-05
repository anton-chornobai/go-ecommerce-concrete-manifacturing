package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/google/uuid"
)

func (p *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	var width, height sql.NullInt64

	query := `SELECT 
		p.id,
		p.price,
		p.title,
		p.type,
		p.color,
		p.description,
		p.status,
		p.stock_quantity,
		p.weight_grams,
		p.rating,
		p.size_width,
		p.size_height,
		p.created_at,
		p.updated_at,
		pi.id,
		pi.url,
		pi.position,
		pi.created_at
	FROM products p
	LEFT JOIN product_image pi ON pi.product_id = p.id
	WHERE p.id=$1
	ORDER BY pi.position ASC
	`
	rows, err := p.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("GetByID: query error: %w", err)
	}
	defer rows.Close()
	found := false

	for rows.Next() {
		var (
			imgID        uuid.NullUUID
			imgURL       sql.NullString
			imgPosition  sql.NullInt16
			imgCreatedAt sql.NullTime
		)

		err := rows.Scan(
			&product.ID,
			&product.Price,
			&product.Title,
			&product.Type,
			&product.Color,
			&product.Description,
			&product.Status,
			&product.StockQuantity,
			&product.Weight,
			&product.Rating,
			&width,
			&height,
			&product.CreatedAt,
			&product.UpdatedAt,
			&imgID,
			&imgURL,
			&imgPosition,
			&imgCreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("GetByID: scan error: %w", err)
		}

		found = true

		if imgID.Valid {
			product.ImageURLs = append(product.ImageURLs, domain.ProductImage{
				ID:        imgID.UUID,
				Key:       imgURL.String,
				Position:  int(imgPosition.Int16),
				CreatedAt: imgCreatedAt.Time,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetByID: rows error: %w", err)
	}

	if !found {
		return nil, fmt.Errorf("product not found")
	}

	if width.Valid && height.Valid {
		product.Size = &domain.Size{
			Width:  int(width.Int64),
			Height: int(height.Int64),
		}
	}

	return &product, nil
}


