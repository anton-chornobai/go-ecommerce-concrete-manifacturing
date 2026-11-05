package application

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
)

func (p *ProductService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.ProductPatchRequest,
	file multipart.File,
	header *multipart.FileHeader,
) error {
	patch := &domain.ProductPatch{
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

	var uploadedKey string

	if file != nil {
		if header == nil {
			return errors.New("file header is required when uploading an image")
		}

		key, err := p.imageStorage.Upload(
			ctx,
			file,
			header.Filename,
			header.Header.Get("Content-Type"),
		)
		if err != nil {
			return fmt.Errorf("upload product image: %w", err)
		}

		uploadedKey = key
		patch.ImageURL = &uploadedKey
	}

	if err := p.repo.Patch(ctx, id, patch); err != nil {
		// The new image was uploaded, but the database update failed.
		// Remove the newly uploaded object to avoid leaving an orphan.
		if uploadedKey != "" {
			cleanupErr := p.cleanupImageKeys([]string{uploadedKey})
			return errors.Join(
				fmt.Errorf("update product %s: %w", id, err),
				cleanupErr,
			)
		}

		return fmt.Errorf("update product %s: %w", id, err)
	}

	return nil
}