package product_handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/product/application"
)

type ProductPostRequest struct {
	Price         int      `schema:"price"`
	Title         string   `schema:"title"`
	Type          string   `schema:"type"`
	Status        string   `schema:"status"`
	Color         *string  `schema:"color"`
	Description   *string  `schema:"description"`
	StockQuantity *int     `schema:"stock_quantity"`
	Weight        *int     `schema:"weight"`
	Rating        *int     `schema:"rating"`
	SizeWidth     *int     `schema:"width"`
	SizeHeight    *int     `schema:"height"`
	File          *os.File `schema:"-"`
}

const imageMaxSize = 10 << 20

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	r.Body = http.MaxBytesReader(w, r.Body, imageMaxSize)
	parsedProduct := new(ProductPostRequest)

	if err := decoder.Decode(parsedProduct, r.Form); err != nil {
		var maxBytesErorr *http.MaxBytesError
		if errors.As(err, maxBytesErorr) {
			http.Error(w, "image should be less than 10MiB", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	productInput := application.CreateProductInput{
		Price:         parsedProduct.Price,
		Title:         parsedProduct.Title,
		Type:          parsedProduct.Type,
		Status:        parsedProduct.Status,
		Color:         parsedProduct.Color,
		Description:   parsedProduct.Description,
		StockQuantity: parsedProduct.StockQuantity,
		Weight:        parsedProduct.Weight,
		Rating:        parsedProduct.Rating,
		SizeWidth:     parsedProduct.SizeWidth,
		SizeHeight:    parsedProduct.SizeHeight,
	}

	var fileInputs []application.FileInput

	for _, header := range r.MultipartForm.File["attachments"] {
		file, err := header.Open()
		if err != nil {
			http.Error(w, "Failed to open uploaded file", http.StatusBadRequest)
            return
		}

		defer file.Close()

		fileInputs = append(fileInputs, application.FileInput{
			Reader: file,
			ContentType: header.Header.Get("Content-Type"),
			Filename: header.Filename,
		})
	}

	product, err := h.ProductService.Create(ctx, productInput, fileInputs)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{ 
		"data": product,
	}); err != nil {
		h.logger.Error("failed to encode product response", "error", err)
		return
	}
}
