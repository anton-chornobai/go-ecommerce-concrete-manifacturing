package application

import (
	"context"
	"errors"
	"fmt"

	"mime/multipart"
	"time"

	"os"
	"strings"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
)

type ProductService struct {
	imageStorage domain.GCSUploader
	repo         domain.Repository
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

func NewProductService(repo domain.Repository, uploader domain.GCSUploader) (*ProductService, error) {
	return &ProductService{repo: repo, imageStorage: uploader}, nil
}

func (p *ProductService) GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error) {
	products, err := p.repo.GetProducts(ctx, limit, status)

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

	return product, nil
}

func (p *ProductService) Add(
    ctx context.Context,
    input *dto.ProductPostRequest,
    file multipart.File,
    header *multipart.FileHeader,
) error {
    if input == nil {
        return errors.New("input data is required")
    }

    var imageURL *string
    if file != nil && header != nil {

        filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)
        url, err := p.imageStorage.Upload(ctx, file, filename)
        if err != nil {
            return fmt.Errorf("upload failed: %w", err)
        }
        imageURL = &url
    }

    var size *domain.Size
    if input.SizeHeight != nil && input.SizeWidth != nil {
        size = &domain.Size{
            Width:  *input.SizeWidth,
            Height: *input.SizeHeight,
        }
    }
	fmt.Printf("%v, %v, %v", input, imageURL, file)
    product, err := domain.NewProduct(
        input.Price,
        input.Title,
        input.Type,
        input.Color,
        input.Status,
        imageURL,
        input.StockQuantity,
        input.Description,
        input.Weight,
        input.Rating,
        size,
    )
    if err != nil {
        return err
    }

    return p.repo.Add(ctx, product)
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

func (p *ProductService) Update(
	ctx context.Context,
	id int,
	req ProductPatchRequest,
	file multipart.File,
	header *multipart.FileHeader,
) error {

	if file != nil && header != nil {
		defer file.Close()

		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)

		url, err := p.imageStorage.Upload(ctx, file, filename)
		if err != nil {
			return fmt.Errorf("upload failed: %w", err)
		}

		req.ImageURL = &url
	}

	update := domain.ProductUpdate{
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

	return p.repo.Update(ctx, id, update)
}