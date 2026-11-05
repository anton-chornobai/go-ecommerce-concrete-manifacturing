package application

import (
	"context"
	"io"
	"log/slog"
	"github.com/google/uuid"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductService struct {
	imageStorage ImageStorage
	repo         domain.Repository
	logger       *slog.Logger
}

type ImageStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	Delete(ctx context.Context, imgURL string) error
}

func NewProductService(repo domain.Repository, storage ImageStorage, logger *slog.Logger) (*ProductService, error) {
	return &ProductService{repo: repo, imageStorage: storage, logger: logger}, nil
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

func (p *ProductService) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := p.repo.GetByID(ctx, id)
	if err != nil {
		p.logger.ErrorContext(ctx, "Не вдалося отримати продукт за id",
			slog.String("помилка", err.Error()),
		)
		return nil, err
	}
	return product, nil
}
