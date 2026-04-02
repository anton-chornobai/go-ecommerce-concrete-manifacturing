package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"strings"

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
			type,
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

func (p *ProductRepository) GetWithLimit(ctx context.Context, limit int) ([]domain.Product, error) {
	var products []domain.Product
	query := `SELECT 
		id,
		price,
		title,
		type,
		image_url,
		color,
		description,
		status,
		stock_quantity,
		weight_grams,
		rating,
		size_width,
		size_height
	FROM products LIMIT $1
	`

	rows, err := p.DB.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		var width, height sql.NullInt64

		err = rows.Scan(
			&product.ID,
			&product.Price,
			&product.Title,
			&product.Type,
			&product.ImageURL,
			&product.Color,
			&product.Description,
			&product.Status,
			&product.StockQuantity,
			&product.Weight,
			&product.Rating,
			&width,
			&height,
		)

		if err != nil {
			return nil, fmt.Errorf("fail scanning product from DB: %w", err)
		}
		if width.Valid && height.Valid {
			product.Size = &domain.Size{
				Width:  int(width.Int64),
				Height: int(height.Int64),
			}
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var product domain.Product
	var width, height sql.NullInt64

	query := `SELECT 
		id,
		price,
		title,
		type,
		image_url,
		color,
		description,
		status,
		stock_quantity,
		weight_grams,
		rating,
		size_width,
		size_height
	FROM products WHERE id=$1;
	`

	row := p.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&product.ID,
		&product.Price,
		&product.Title,
		&product.Type,
		&product.ImageURL,
		&product.Color,
		&product.Description,
		&product.Status,
		&product.StockQuantity,
		&product.Weight,
		&product.Rating,
		&width,
		&height,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("not found row")
		}
		return nil, fmt.Errorf("GetByID: scanning error: %w", err)
	}

	if width.Valid && height.Valid {
		product.Size = &domain.Size{
			Width:  int(width.Int64),
			Height: int(height.Int64),
		}
	} else {
		product.Size = nil
	}

	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, id int, req domain.ProductUpdate) error {
	var sb strings.Builder
	args := []any{}
	i := 1

	sb.WriteString("UPDATE products SET ")

	fields := []struct {
		column string
		value  any
	}{
		{"price", req.Price},
		{"title", req.Title},
		{"type", req.ProductType},
		{"image_url", req.ImageURL},
		{"color", req.Color},
		{"status", req.Status},
		{"description", req.Description},
		{"stock_quantity", req.StockQuantity},
		{"weight_grams", req.WeightGrams},
		{"rating", req.Rating},
		{"size_width", req.SizeWidth},
		{"size_height", req.SizeHeight},
	}

	for _, f := range fields {
		if f.value != nil {
			fmt.Fprintf(&sb, "%s=$%d,", f.column, i)
			args = append(args, f.value)
			i++
		}
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := strings.TrimSuffix(sb.String(), ",")
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	res, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (p *ProductRepository) DeleteByID(ctx context.Context, id int) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE order_item SET product_id = NULL WHERE product_id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to nullify product references: %w", err)
	}

	res, err := tx.ExecContext(ctx,
		`DELETE FROM products WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
