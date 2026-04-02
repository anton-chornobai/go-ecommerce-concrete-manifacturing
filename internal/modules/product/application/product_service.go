package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductService struct {
	repo domain.Repository
	log  *slog.Logger
}

type ProductPatchRequest struct {
	Price         *int                  `json:"price,omitempty"`
	Title         *string               `json:"title,omitempty"`
	ProductType   *string               `json:"type,omitempty"`
	ImageURL      *string               `json:"image_url,omitempty"`
	Color         *string               `json:"color,omitempty"`
	Status        *domain.ProductStatus `json:"status,omitempty"`
	Description   *string               `json:"description,omitempty"`
	StockQuantity *int                  `json:"stock_quantity,omitempty"`
	WeightGrams   *int                  `json:"weight,omitempty"`
	Rating        *int                  `json:"rating,omitempty"`
	SizeWidth     *int                  `json:"size_width,omitempty"`
	SizeHeight    *int                  `json:"size_height,omitempty"`
}

func NewProductService(repo domain.Repository) (*ProductService, error) {
	return &ProductService{repo: repo}, nil
}

func (p *ProductService) GetWithLimit(ctx context.Context, limit int) ([]domain.Product, error) {
	products, err := p.repo.GetWithLimit(ctx, limit)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (p *ProductService) GetById(ctx context.Context, id int) (*domain.Product, error) {
	product, err := p.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return  product, nil
}

func (p *ProductService) Add(ctx context.Context, input domain.Product) error {
	product, err := domain.NewProduct(
		input.Price,
		input.Title,
		input.Type,
		input.Color,
		domain.ProductStatus(input.Status),
		input.ImageURL,
		input.StockQuantity,
		input.Description,
		input.Weight,
		input.Rating,
		input.Size,
	)

	if err != nil {
		return err
	}

	err = p.repo.Add(ctx, product)

	if err != nil {
		return err
	}

	return nil
}

func (p *ProductService) DeleteByID(ctx context.Context, id int) error {
	product, err := p.repo.GetByID(ctx, id)

	if err != nil {
		return fmt.Errorf("failed to get product form db %w", err)
	}

	err = p.repo.DeleteByID(ctx, id)

	if err != nil {
		return err
	}

	if product.ImageURL != nil {
		filePath := strings.TrimPrefix(*product.ImageURL, "/")
		if err := os.Remove(filePath); err != nil && errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("could not delete image: %w", err)
		}
	}

	return nil
}

func (p *ProductService) Update(ctx context.Context, id int, req ProductPatchRequest) error {

	product := domain.ProductUpdate{
		Price:         req.Price,
		Title:         req.Title,
		ProductType:   req.ProductType,
		ImageURL:      req.ImageURL,
		Color:         req.Color,
		Status:        req.Status,
		Description:   req.Description,
		StockQuantity: req.StockQuantity,
		WeightGrams:   req.WeightGrams,
		Rating:        req.Rating,
		SizeWidth:     req.SizeWidth,
		SizeHeight:    req.SizeHeight,
	}
	err := p.repo.Update(ctx, id, product)

	if err != nil {
		return err
	}

	return nil
}
