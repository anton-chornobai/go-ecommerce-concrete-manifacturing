package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	logger       *slog.Logger
}

func NewProductService(repo domain.Repository, uploader domain.GCSUploader, logger *slog.Logger) (*ProductService, error) {
	return &ProductService{repo: repo, imageStorage: uploader, logger: logger}, nil
}
func (p *ProductService) GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error) {
	products, err := p.repo.GetProducts(ctx, limit, status)
	if err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося отримати список продуктів",
			slog.Int("limit", limit),
			slog.Any("status", status),
			slog.String("помилка", err.Error()),
		)
		return nil, err
	}

	p.logger.InfoContext(ctx, "Список продуктів отримано",
		slog.Int("кількість", len(products)),
	)
	return products, nil
}

func (p *ProductService) GetById(ctx context.Context, id int) (*domain.Product, error) {
	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося отримати продукт за id",
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
		return nil, err
	}

	p.logger.InfoContext(ctx, "Продукт отримано",
		slog.Int("id", id),
	)
	return product, nil
}

func (p *ProductService) Add(
	ctx context.Context,
	input *dto.ProductPostRequest,
	file multipart.File,
	header *multipart.FileHeader,
) error {
	if input == nil {
		p.logger.ErrorContext(ctx, "Вхідні дані відсутні")
		return errors.New("input data is required")
	}

	var imageURL *string
	if file != nil && header != nil {
		filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)
		url, err := p.imageStorage.Upload(ctx, file, filename)
		if err != nil {
			p.logger.ErrorContext(ctx, "Не вдалося завантажити зображення",
				slog.String("файл", filename),
				slog.String("помилка", err.Error()),
			)
			return fmt.Errorf("upload failed: %w", err)
		}
		p.logger.InfoContext(ctx, "Зображення завантажено",
			slog.String("url", url),
		)
		imageURL = &url
	}

	var size *domain.Size
	if input.SizeHeight != nil && input.SizeWidth != nil {
		size = &domain.Size{
			Width:  *input.SizeWidth,
			Height: *input.SizeHeight,
		}
	}

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
		p.logger.ErrorContext(ctx, "Не вдалося створити продукт — помилка валідації",
			slog.String("назва", input.Title),
			slog.String("помилка", err.Error()),
		)
		return err
	}

	if err := p.repo.Add(ctx, product); err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося зберегти продукт у базі даних",
			slog.String("назва", input.Title),
			slog.String("помилка", err.Error()),
		)
		return err
	}

	p.logger.InfoContext(ctx, "Продукт успішно створено",
		slog.String("назва", input.Title),
	)
	return nil
}

func (p *ProductService) DeleteByID(ctx context.Context, id int) error {
	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося знайти продукт для видалення",
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
		return fmt.Errorf("failed to get product form db %w", err)
	}

	if err := p.repo.DeleteByID(ctx, id); err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося видалити продукт з бази даних",
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
		return err
	}

	if product.ImageURL != nil {
		filePath := strings.TrimPrefix(*product.ImageURL, "/")
		if err := os.Remove(filePath); err != nil && errors.Is(err, os.ErrNotExist) {
			p.logger.ErrorContext(ctx, "Не вдалося видалити зображення продукту",
				slog.Int("id", id),
				slog.String("шлях", filePath),
				slog.String("помилка", err.Error()),
			)
			return fmt.Errorf("could not delete image: %w", err)
		}
	}

	p.logger.InfoContext(ctx, "Продукт успішно видалено",
		slog.Int("id", id),
	)
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
			p.logger.ErrorContext(ctx, "Не вдалося завантажити зображення для оновлення",
				slog.Int("id", id),
				slog.String("файл", uniqueName),
				slog.String("помилка", err.Error()),
			)
			return err
		}
		p.logger.InfoContext(ctx, "Зображення для продукту оновлено",
			slog.Int("id", id),
			slog.String("url", url),
		)
		patch.ImageURL = &url
	}

	if err := patch.Validate(); err != nil {
		p.logger.WarnContext(ctx, "Невалідні дані для оновлення продукту",
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
		return err
	}

	if err := p.repo.Patch(ctx, id, &patch); err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося оновити продукт у базі даних",
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
		return err
	}

	p.logger.InfoContext(ctx, "Продукт успішно оновлено",
		slog.Int("id", id),
	)
	return nil
}
