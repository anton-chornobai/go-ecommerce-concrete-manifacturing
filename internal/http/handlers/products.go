package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
	"github.com/gorilla/schema"
)

const (
	// To limit body (with file) size in Update and Add
	maxBytesBodyLimit = 10 << 20
)

type ProductHandler struct {
	ProductService application.ProductService
	logger         *slog.Logger
}

func NewProductsHandler(productService application.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{ProductService: productService, logger: logger}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var status *domain.ProductStatus

	queryParams := r.URL.Query()
	statusVal := queryParams.Get("status")

	if statusVal == "" ||
		(statusVal != string(domain.ProductArchived) &&
			statusVal != string(domain.ProductDisplayed)) {
		status = nil
	} else {
		s := domain.ProductStatus(statusVal)
		status = &s
	}

	products, err := h.ProductService.GetProducts(ctx, 20, status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string][]domain.Product{
		"data": products,
	}); err != nil {
		http.Error(w, "couldnt send data", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	idStr := pathParts[3]
	fmt.Println(idStr)

	ingtegerID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	product, err := h.ProductService.GetById(ctx, ingtegerID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "no such product", http.StatusBadRequest)
			return
		}
		fmt.Println(err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"data": product,
	}); err != nil {
		http.Error(w, "couldnt encode response", http.StatusInternalServerError)
		return

	}
}

var decoder = schema.NewDecoder()

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	r.Body = http.MaxBytesReader(w, r.Body, maxBytesBodyLimit)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, fmt.Sprintf("Тіло запиту занадто велике (більше %v)", maxBytesBodyLimit), http.StatusRequestEntityTooLarge)
			return
		}
	}
	parsedForm := new(dto.ProductPostRequest)

	if err := decoder.Decode(parsedForm, r.PostForm); err != nil {
		http.Error(w, "Помилка декодування"+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Помилка отримання файлу", http.StatusBadRequest)
			return
		}

		file = nil
		header = nil
	} else {
		defer file.Close()
	}

	if err := h.ProductService.Add(ctx, parsedForm, file, header); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Продукт створений",
	})

	if err != nil {
		h.logger.Warn("ProductHandler.Add Failed to encode JSON response", "err:", err.Error())
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	//
}

func (h *ProductHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(idStr[len(idStr)-1])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.ProductService.DeleteByID(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "deleted",
	})
}
