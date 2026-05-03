package dto

import (
	"mime/multipart"
	"os"

	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

type ProductPostRequest struct {
	Price         int                  `schema:"price"`
	Title         string               `schema:"title"`
	Type          string               `schema:"type"`
	Status        domain.ProductStatus `schema:"status"`
	Color         *string              `schema:"color"`
	Description   *string              `schema:"description"`
	StockQuantity *int                 `schema:"stock_quantity"`
	Weight        *int                 `schema:"weight"`
	Rating        *int                 `schema:"rating"`
	SizeWidth     *int                 `schema:"width"`
	SizeHeight    *int                 `schema:"height"`
	File          *os.File             `schema:"-"`
}

type ProductPatchRequest struct {
	Price           *int                  `schema:"price"`
	Title           *string               `schema:"title"`
	Type            *string               `schema:"type"`
	Status          *domain.ProductStatus `schema:"status"`
	Color           *string               `schema:"color"`
	Description     *string               `schema:"description"`
	StockQuantity   *int                  `schema:"stock_quantity"`
	Weight          *int                  `schema:"weight"`
	Rating          *int                  `schema:"rating"`
	SizeWidth       *int                  `schema:"width"`
	SizeHeight      *int                  `schema:"height"`
	ImageFileHeader *multipart.FileHeader `schema:"-"`
}
