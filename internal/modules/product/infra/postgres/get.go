package postgres

import (
	"context"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

func (p *ProductRepository) GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error) {
	_ = `
		SELECT
			p.id, p.price, p.title, p.type, p.status,
			p.color, p.description, p.stock_quantity,
			p.weight_grams, p.rating, p.size_width, p.size_height,
			p.created_at, p.updated_at,
			pi.id, pi.url, pi.position, pi.created_at
		FROM products p
		LEFT JOIN product_image pi ON pi.product_id = p.id`


	return  nil, nil
}