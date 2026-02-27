package domain

import "errors"

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Product struct {
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
	Size          *Size   `json:"size,omitempty"`
}

func NewProduct(
	price int,
	title string,
	productType string,
	imageURL string,
	color string,
	stockQuantity *int,
	description *string,
	weight *int,
	rating *int,
	size *Size,
) (*Product, error) {

	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	if len(title) < 2 {
		return nil, errors.New("title must be at least 2 characters")
	}

	if productType == "" {
		return nil, errors.New("product type is required")
	}

	if stockQuantity != nil && *stockQuantity < 0 {
		return nil, errors.New("stock quantity cannot be negative")
	}

	if weight != nil && *weight < 0 {
		return nil, errors.New("weight cannot be negative")
	}

	return &Product{
		Price:         price,
		Title:         title,
		Type:          productType,
		ImageURL:      imageURL,
		Color:         color,
		StockQuantity: stockQuantity,
		Description:   description,
		Weight:        weight,
		Rating:        rating,
		Size:          size,
	}, nil
}
