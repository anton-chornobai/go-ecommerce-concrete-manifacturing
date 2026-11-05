package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/google/uuid"
)

func (r *ProductRepository) Patch(ctx context.Context, id uuid.UUID, p *domain.ProductPatch) error {
	setParts := []string{}
	args := []any{}
	i := 1

	add := func(column string, value any) {
		setParts = append(setParts, fmt.Sprintf("%s=$%d", column, i))
		args = append(args, value)
		i++
	}

	if p.Price != nil {
		add("price", *p.Price)
	}
	if p.Title != nil {
		add("title", *p.Title)
	}
	if p.Type != nil {
		add("type", *p.Type)
	}
	if p.Status != nil {
		add("status", *p.Status)
	}
	if p.ImageURL != nil {
		add("image_url", *p.ImageURL)
	}
	if p.Color != nil {
		add("color", *p.Color)
	}
	if p.Description != nil {
		add("description", *p.Description)
	}
	if p.StockQuantity != nil {
		add("stock_quantity", *p.StockQuantity)
	}
	if p.Weight != nil {
		add("weight_grams", *p.Weight)
	}
	if p.Rating != nil {
		add("rating", *p.Rating)
	}
	if p.SizeWidth != nil {
		add("size_width", *p.SizeWidth)
	}
	if p.SizeHeight != nil {
		add("size_height", *p.SizeHeight)
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE products SET %s WHERE id=$%d",
		strings.Join(setParts, ", "),
		i,
	)
	args = append(args, id)

	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to patch product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}
