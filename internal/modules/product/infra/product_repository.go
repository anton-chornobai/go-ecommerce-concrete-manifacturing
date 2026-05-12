package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"strings"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProductRepository struct {
	DB *sql.DB
}

func (p *ProductRepository) Add(ctx context.Context, product *domain.Product) error {
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
			image.URL,
			i,
		)
		if err != nil {
			return fmt.Errorf("не вдалося зберегти зображення %s: %w", image.URL, err)
		}
	}

	return nil
}
func (p *ProductRepository) GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error) {
	query := `
		SELECT
			p.id, p.price, p.title, p.type, p.status,
			p.color, p.description, p.stock_quantity,
			p.weight_grams, p.rating, p.size_width, p.size_height,
			p.created_at, p.updated_at,
			pi.id, pi.url, pi.position, pi.created_at
		FROM products p
		LEFT JOIN product_image pi ON pi.product_id = p.id`

	args := []any{}
	argPos := 1

	if status != nil {
		query += fmt.Sprintf(" WHERE p.status=$%d", argPos)
		args = append(args, *status)
		argPos++
	}

	query += " ORDER BY p.id, pi.position"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argPos)
		args = append(args, limit)
		argPos++
	}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не вдалося виконати запит: %w", err)
	}
	defer rows.Close()

	productMap := map[int]*domain.Product{}
	productOrder := []int{}

	for rows.Next() {
		var (
			sizeWidth, sizeHeight          *int
			imageID                        *uuid.UUID
			imageURL                       *string
			imagePosition                  *int
			imageCreatedAt                 *time.Time
		)

		// temp product to scan into each row
		var row domain.Product

		err := rows.Scan(
			&row.ID,
			&row.Price,
			&row.Title,
			&row.Type,
			&row.Status,
			&row.Color,
			&row.Description,
			&row.StockQuantity,
			&row.Weight,
			&row.Rating,
			&sizeWidth,
			&sizeHeight,
			&row.CreatedAt,
			&row.UpdatedAt,
			&imageID,
			&imageURL,
			&imagePosition,
			&imageCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не вдалося зчитати рядок: %w", err)
		}

		// first time seeing this product id — add to map
		if _, exists := productMap[row.ID]; !exists {
			if sizeWidth != nil && sizeHeight != nil {
				row.Size = &domain.Size{
					Width:  *sizeWidth,
					Height: *sizeHeight,
				}
			}
			productMap[row.ID] = &row
			productOrder = append(productOrder, row.ID)
		}

		// append image if this row has one
		if imageID != nil {
			productMap[row.ID].ImageURLs = append(productMap[row.ID].ImageURLs, domain.ProductImage{
				ID:        *imageID,
				ProductID: row.ID,
				URL:       *imageURL,
				Position:  *imagePosition,
				CreatedAt: *imageCreatedAt,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("помилка ітерації рядків: %w", err)
	}

	// preserve order using the slice
	products := make([]domain.Product, 0, len(productOrder))
	for _, id := range productOrder {
		products = append(products, *productMap[id])
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
		&product.ImageURLs,
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

func (r *ProductRepository) Patch(ctx context.Context, id int, p *domain.ProductPatch) error {
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
