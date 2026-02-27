package handlers

import (
	"context"
	"encoding/json"

	"net/http"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductHandler struct {
	ProductService application.ProductService
}

type ProductRequest struct {
	ID            int     `json:"id"`
	Price         int     `json:"price"`
	Title         string  `json:"title"`
	Type          string  `json:"type"`
	ImageURL      string  `json:"imageUrl"`
	Color         string  `json:"color"`
	Description   *string `json:"description,omitempty"`
	StockQuantity *int    `json:"stockQuantity,omitempty"`
	Weight        *int    `json:"weight,omitempty"`
	Rating        *int    `json:"rating,omitempty"`
	Size          *struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size,omitempty"`
}

func NewProductsHandler(productService application.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
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
	if req.Size != nil {
		size = &domain.Size{
			Width:  req.Size.Width,
			Height: req.Size.Height,
		}
	}

	product := domain.Product{
		Price:         req.Price,
		Title:         req.Title,
		Type:          req.Type,
		ImageURL:      req.ImageURL,
		Color:         req.Color,
		Description:   req.Description,
		StockQuantity: req.StockQuantity,
		Weight:        req.Weight,
		Rating:        req.Rating,
		Size:          size,
	}

	productID, err := h.ProductService.Add(ctx, product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"created_id": productID,
		"title": product.Title,
	})
}