package product_handler

import (
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"github.com/google/uuid"
	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
	"github.com/gorilla/schema"
)

const (
	maxBytesBodyLimit = 10 << 20 // 10MB
)

var decoder = schema.NewDecoder()

type ProductHandler struct {
	ProductService ProductService
	logger         *slog.Logger
}

type ProductService interface {
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetProducts(ctx context.Context, limit int, status *domain.ProductStatus) ([]domain.Product, error)
	Update(
		ctx context.Context,
		id uuid.UUID, req dto.ProductPatchRequest,
		file multipart.File,
		header *multipart.FileHeader,
	) error
	Create(
		ctx context.Context,
		input application.CreateProductInput,
		fileInput []application.FileInput,
	) (domain.Product, error)
}

func NewProductsHandler(productService ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{ProductService: productService, logger: logger}
}

func (h *ProductHandler) handleServiceError(w http.ResponseWriter, r *http.Request, err error) {

	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		h.logger.WarnContext(r.Context(), "Продукт не знайдено",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "handleServiceError")),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrTitleAlreadyExists):
		http.Error(w, "назва продукту уже існує", http.StatusInternalServerError)

	case errors.Is(err, domain.ErrInvalidPrice),
		errors.Is(err, domain.ErrTitleTooShort),
		errors.Is(err, domain.ErrTitleTooLong),
		errors.Is(err, domain.ErrTypeRequired),
		errors.Is(err, domain.ErrInvalidStatus),
		errors.Is(err, domain.ErrNegativeStock),
		errors.Is(err, domain.ErrNegativeWeight),
		errors.Is(err, domain.ErrNegativeHeight),
		errors.Is(err, domain.ErrNegativeWidth):
		h.logger.WarnContext(r.Context(), "Невалідні дані від клієнта",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "handleServiceError")),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		h.logger.ErrorContext(r.Context(), "Внутрішня помилка сервісу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "handleServiceError")),
			slog.String("шлях", r.URL.Path),
			slog.String("метод", r.Method),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Щось пішло не так", http.StatusInternalServerError)
	}
}
