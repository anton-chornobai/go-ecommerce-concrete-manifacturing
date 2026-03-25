package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"net/http"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductHandler struct {
	ProductService application.ProductService
}

func NewProductsHandler(productService application.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

type ProductRequest struct {
	ID            int     `json:"id"`
	Price         int     `json:"price"`
	Title         string  `json:"title"`
	Type          string  `json:"type"`
	Color         string  `json:"color"`
	Status        *string `json:"status,omitempty"`
	ImageURL      *string `json:"image_url"`
	Description   *string `json:"description,omitempty"`
	StockQuantity *int    `json:"stock_quantity,omitempty"`
	Weight        *int    `json:"weight,omitempty"`
	Rating        *int    `json:"rating,omitempty"`
	Width         *int    `json:"width,omitempty"`
	Height        *int    `json:"height,omitempty"`
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	products, err := h.ProductService.GetWithLimit(ctx, 20)

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

func (h *ProductHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var size *domain.Size
	if req.Width != nil && req.Height != nil {
		size = &domain.Size{
			Width:  *req.Width,
			Height: *req.Height,
		}
	}

	var productSatusSafe domain.ProductStatus

	if req.Status != nil {
		productSatusSafe = domain.ProductStatus(*req.Status)
	}

	product := domain.Product{
		Price:         req.Price,
		Title:         req.Title,
		Type:          req.Type,
		ImageURL:      req.ImageURL,
		Status:        productSatusSafe,
		Color:         req.Color,
		Description:   req.Description,
		StockQuantity: req.StockQuantity,
		Weight:        req.Weight,
		Rating:        req.Rating,
		Size:          size,
	}

	err := h.ProductService.Add(ctx, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": product.Title,
	})
}

func (h *ProductHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.ProductService.DeleteByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Product deleted",
	})
	if err != nil {
		http.Error(w, "couldn't write response about deletion", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var req application.ProductPatchRequest

	path := r.URL.Path

	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	idStr := parts[3]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	err = h.ProductService.Update(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
	})

	if err != nil {
		http.Error(w, "failed to send response message", http.StatusInternalServerError)
		return
	}
}
