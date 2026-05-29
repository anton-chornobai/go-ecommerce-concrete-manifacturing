package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
	"github.com/gorilla/schema"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

const (
	maxBytesBodyLimit = 10 << 20 // 10MB
)

var decoder = schema.NewDecoder()

type ProductHandler struct {
	ProductService application.ProductService
	logger         *slog.Logger
}

func NewProductsHandler(productService application.ProductService, logger *slog.Logger) *ProductHandler {
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
	case  errors.Is(err, domain.ErrTitleAlreadyExists):
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

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	h.logger.InfoContext(ctx, "Запит на отримання списку продуктів",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
	)

	var status *domain.ProductStatus
	statusVal := r.URL.Query().Get("status")
	if statusVal == string(domain.ProductArchived) || statusVal == string(domain.ProductDisplayed) {
		s := domain.ProductStatus(statusVal)
		status = &s
	}

	products, err := h.ProductService.GetProducts(ctx, 20, status)
	if err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося отримати список продуктів",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Щось пішло не так", http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(ctx, "Список продуктів успішно отримано",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
		slog.Int("кількість", len(products)),
	)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]domain.Product{"data": products}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
			slog.String("помилка", err.Error()),
		)
	}
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strID)
	if err != nil || id <= 0 {
		h.logger.WarnContext(ctx, "Невалідний id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
			slog.String("отримано", strID),
		)
		http.Error(w, "id має бути цілим додатнім числом", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на отримання продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		slog.Int("id", id),
	)

	product, err := h.ProductService.GetById(ctx, id)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно отримано",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		slog.Int("id", id),
	)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]any{"data": product}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
	}
}

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	r.Body = http.MaxBytesReader(w, r.Body, maxBytesBodyLimit)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.logger.WarnContext(ctx, "Тіло запиту занадто велике",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Add")),
				slog.Int64("ліміт_байт", maxBytesErr.Limit),
			)
			http.Error(w, fmt.Sprintf("Тіло запиту занадто велике (максимум %dmb)", maxBytesBodyLimit>>20), http.StatusRequestEntityTooLarge)
		} else {
			h.logger.WarnContext(ctx, "Неправильний формат тіла запиту",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Add")),
				slog.String("помилка", err.Error()),
			)
			http.Error(w, "Неправильний формат запиту", http.StatusBadRequest)
		}
		return
	}

	parsedForm := new(dto.ProductPostRequest)
	if err := decoder.Decode(parsedForm, r.PostForm); err != nil {
		h.logger.WarnContext(ctx, "Помилка декодування форми",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Add")),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Помилка декодування: "+err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["attachments"]
	fmt.Println(len(files))
	if len(files) < 1 {
		http.Error(w, "Потрібно мінімум одне фото для продукту щоб його створити", http.StatusBadRequest)
		return
	}

	if err := h.ProductService.Add(ctx, parsedForm, files); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно створено",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Add")),
		slog.String("назва", parsedForm.Title),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Продукт створений"}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь після створення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Add")),
			slog.String("помилка", err.Error()),
		)
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	intID, err := strconv.Atoi(strID)
	if err != nil || intID <= 0 {
		h.logger.WarnContext(ctx, "Невалідний id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.String("отримано", strID),
		)
		http.Error(w, "id має бути цілим додатнім числом", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на оновлення продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		slog.Int("id", intID),
	)

	r.Body = http.MaxBytesReader(w, r.Body, maxBytesBodyLimit)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.logger.WarnContext(ctx, "Тіло запиту занадто велике при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.Int("id", intID),
				slog.Int64("ліміт_байт", maxBytesErr.Limit),
			)
			http.Error(w, fmt.Sprintf("Тіло запиту занадто велике (максимум %dmb)", maxBytesBodyLimit>>20), http.StatusRequestEntityTooLarge)
		} else {
			h.logger.WarnContext(ctx, "Неправильний формат тіла запиту при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.Int("id", intID),
				slog.String("помилка", err.Error()),
			)
			http.Error(w, "Неправильний формат запиту", http.StatusBadRequest)
		}
		return
	}

	parseForm := new(dto.ProductPatchRequest)
	if err := decoder.Decode(parseForm, r.PostForm); err != nil {
		h.logger.WarnContext(ctx, "Помилка декодування форми оновлення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.Int("id", intID),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Помилка декодування: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			file = nil
			header = nil
		case errors.Is(err, multipart.ErrMessageTooLarge):
			h.logger.WarnContext(ctx, "Файл завеликий при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.Int("id", intID),
			)
			http.Error(w, "Файл завеликий", http.StatusRequestEntityTooLarge)
			return
		default:
			h.logger.WarnContext(ctx, "Помилка отримання файлу при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.Int("id", intID),
				slog.String("помилка", err.Error()),
			)
			http.Error(w, "Неправильний запит", http.StatusBadRequest)
			return
		}
	}
	if file != nil {
		defer file.Close()
	}

	if err := h.ProductService.Update(ctx, intID, *parseForm, file, header); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно оновлено",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		slog.Int("id", intID),
	)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Продукт з id %d успішно оновлено", intID),
	}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь після оновлення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.Int("id", intID),
			slog.String("помилка", err.Error()),
		)
	}
}

func (h *ProductHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strID)
	if err != nil || id <= 0 {
		h.logger.WarnContext(ctx, "Невалідний id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
			slog.String("отримано", strID),
		)
		http.Error(w, "id має бути цілим додатнім числом", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на видалення продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		slog.Int("id", id),
	)

	if err := h.ProductService.DeleteByID(ctx, id); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно видалено",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		slog.Int("id", id),
	)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Продукт з id %d успішно видалено", id),
	}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь після видалення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
			slog.Int("id", id),
			slog.String("помилка", err.Error()),
		)
	}
}
