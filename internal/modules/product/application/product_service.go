package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

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
	req dto.ProductPatchRequest,
	file multipart.File,
	header *multipart.FileHeader,
) error {
	patch := domain.ProductPatch{
		Price:         req.Price,
		Title:         req.Title,
		Type:          req.Type,
		Status:        req.Status,
		Color:         req.Color,
		Description:   req.Description,
		StockQuantity: req.StockQuantity,
		Weight:        req.Weight,
		Rating:        req.Rating,
		SizeWidth:     req.SizeWidth,
		SizeHeight:    req.SizeHeight,
	}

	if file != nil && header != nil {
		ext := filepath.Ext(header.Filename)
		nameOnly := strings.TrimSuffix(header.Filename, ext)
		uniqueName := fmt.Sprintf("%s_%d%s", nameOnly, time.Now().Unix(), ext)

		url, err := p.imageStorage.Upload(ctx, file, uniqueName)
		if err != nil {
			return err
		}
		patch.ImageURL = &url
	}

	if err := patch.Validate(); err != nil {
		return err
	}

	return p.repo.Patch(ctx, id, &patch)
}