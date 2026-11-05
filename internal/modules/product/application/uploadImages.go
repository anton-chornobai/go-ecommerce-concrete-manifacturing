package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

func (p *ProductService) uploadImages(
	ctx context.Context,
	files []FileInput,
) ([]domain.ProductImage, error) {
	images := make([]domain.ProductImage, 0, len(files))
	uploadedKeys := make([]string, 0, len(files))

	cleanup := func(cause error) error {
		cleanupErr := p.cleanupImageKeys(uploadedKeys)
		return errors.Join(cause, cleanupErr)
	}

	for i, file := range files {
		key, err := p.imageStorage.Upload(
			ctx,
			file.Reader,
			file.Filename,
			file.ContentType,
		)

		if err != nil {
			return nil, cleanup(
				fmt.Errorf("upload image %d: %w", i, err),
			)
		}

		uploadedKeys = append(uploadedKeys, key)

		image, err := domain.NewProductImage(i, key)
		if err != nil {
			return nil, cleanup(
				fmt.Errorf("create product image %d: %w", i, err),
			)
		}

		images = append(images, *image)
	}

	return images, nil
}
