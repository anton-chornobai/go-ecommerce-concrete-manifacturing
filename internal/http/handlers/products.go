package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		if errors.Is(err, http.ErrContentLength) {
			http.Error(w, "file too large, max 10 MB", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var imageURL *string

	file, header, err := r.FormFile("image_url")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "failed to read image: "+err.Error(), http.StatusBadRequest)
			return
		}
		//  no file provided, then imageURL stays nil
	} else {
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		allowed := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
			"image/webp": true,
		}
		if !allowed[contentType] {
			http.Error(w, "only jpeg/png/webp allowed", http.StatusBadRequest)
			return
		}

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", filename)

		if err := os.MkdirAll("uploads", 0755); err != nil {
			http.Error(w, "could not create uploads dir", http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "could not save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "could not write file", http.StatusInternalServerError)
			return
		}

		url := "/uploads/" + filename
		imageURL = &url
	}

	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		http.Error(w, "invalid price", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	productType := r.FormValue("type")

	var status *string
	if s := r.FormValue("status"); s != "" {
		status = &s
	}

	var description *string
	if d := r.FormValue("description"); d != "" {
		description = &d
	}

	var stockQuantity *int
	if sq := r.FormValue("stock_quantity"); sq != "" {
		v, err := strconv.Atoi(sq)
		if err != nil {
			http.Error(w, "invalid stock_quantity", http.StatusBadRequest)
			return
		}
		stockQuantity = &v
	}

	var weight *int
	if wv := r.FormValue("weight"); wv != "" {
		v, err := strconv.Atoi(wv)
		if err != nil {
			http.Error(w, "invalid weight", http.StatusBadRequest)
			return
		}
		weight = &v
	}

	var rating *int
	if rv := r.FormValue("rating"); rv != "" {
		v, err := strconv.Atoi(rv)
		if err != nil {
			http.Error(w, "invalid rating", http.StatusBadRequest)
			return
		}
		rating = &v
	}

	var width, height *int
	if wv := r.FormValue("width"); wv != "" {
		v, err := strconv.Atoi(wv)
		if err != nil {
			http.Error(w, "invalid width", http.StatusBadRequest)
			return
		}
		width = &v
	}
	if hv := r.FormValue("height"); hv != "" {
		v, err := strconv.Atoi(hv)
		if err != nil {
			http.Error(w, "invalid height", http.StatusBadRequest)
			return
		}
		height = &v
	}
	var color *string
	if cl := r.FormValue("color"); cl != "" {
		color = &cl
	}

	var size *domain.Size
	if width != nil && height != nil {
		size = &domain.Size{Width: *width, Height: *height}
	}

	var productStatusSafe domain.ProductStatus
	if status != nil {
		productStatusSafe = domain.ProductStatus(*status)
	}

	product := domain.Product{
		Price:         price,
		Title:         title,
		Type:          productType,
		Color:         color,
		ImageURL:      imageURL,
		Status:        productStatusSafe,
		Description:   description,
		StockQuantity: stockQuantity,
		Weight:        weight,
		Rating:        rating,
		Size:          size,
	}

	if err := h.ProductService.Add(ctx, product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"message": "product created"})
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

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	const maxSizeMB = 10
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSizeMB)<<20)

	if err := r.ParseMultipartForm(int64(maxSizeMB) << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	var req application.ProductPatchRequest

	parseInt := func(key string) *int {
		if v := r.FormValue(key); v != "" {
			if i, err := strconv.Atoi(v); err == nil {
				return &i
			}
		}
		return nil
	}

	parseString := func(key string) *string {
		if v := r.FormValue(key); v != "" {
			return &v
		}
		return nil
	}

	req.Price = parseInt("price")
	req.Title = parseString("title")
	req.ProductType = parseString("type")
	req.Color = parseString("color")
	req.Description = parseString("description")

	req.StockQuantity = parseInt("stock_quantity")
	req.WeightGrams = parseInt("weight")
	req.Rating = parseInt("rating")
	req.SizeWidth = parseInt("size_width")
	req.SizeHeight = parseInt("size_height")

	if s := r.FormValue("status"); s != "" {
		status := domain.ProductStatus(s)
		req.Status = &status
	}

	file, header, err := r.FormFile("image_url")
	if err == nil {
		defer file.Close()

		filename := header.Filename
		req.ImageURL = &filename
	}

	err = h.ProductService.Update(ctx, id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "Success",
	})
}
