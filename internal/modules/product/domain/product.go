package domain

import (
	"errors"
	"unicode/utf8"
)

const (
	MaxTitleLength = 30
	MinTitleLength = 2
)

var (
	ErrInvalidPrice   = errors.New("ціна повинна бути більшою за 0")
	ErrTitleTooShort  = errors.New("назва повинна містити принаймні 2 символи")
	ErrTitleTooLong   = errors.New("назва занадто довга (максимум 30 символів)")
	ErrTypeRequired   = errors.New("тип продукту є обов'язковим")
	ErrInvalidStatus  = errors.New("недопустимий статус продукту")
	ErrNegativeStock  = errors.New("кількість на складі не може бути від'ємною")
	ErrNegativeWeight = errors.New("вага не може бути від'ємною")
)

type ProductStatus string

const (
	ProductArchived  ProductStatus = "archived"
	ProductDisplayed ProductStatus = "displayed"
)

func (s ProductStatus) IsValid() bool {
	switch s {
	case ProductArchived, ProductDisplayed:
		return true
	default:
		return false
	}
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Product struct {
	ID            int           `json:"id"`
	Price         int           `json:"price"`
	Title         string        `json:"title"`
	Type          string        `json:"type"`
	Color         *string       `json:"color,omitempty"`
	Status        ProductStatus `json:"status"`
	ImageURL      *string       `json:"image_url,omitempty"`
	Description   *string       `json:"description,omitempty"`
	StockQuantity *int          `json:"stock_quantity,omitempty"`
	Weight        *int          `json:"weight,omitempty"`
	Rating        *int          `json:"rating,omitempty"`
	Size          *Size         `json:"size,omitempty"`
}

type ProductUpdate struct {
	Price         *int
	Title         *string
	ProductType   *string
	ImageURL      *string
	Color         *string
	Status        *ProductStatus
	Description   *string
	StockQuantity *int
	WeightGrams   *int
	Rating        *int
	SizeWidth     *int
	SizeHeight    *int
}

func NewProduct(
	price int,
	title string,
	productType string,
	color *string,
	status ProductStatus,
	imageURL *string,
	stockQuantity *int,
	description *string,
	weight *int,
	rating *int,
	size *Size,
) (*Product, error) {
	if price <= 0 {
		return nil, ErrInvalidPrice
	}

	titleLen := utf8.RuneCountInString(title)
	if titleLen < MinTitleLength {
		return nil, ErrTitleTooShort
	}
	if titleLen > MaxTitleLength {
		return nil, ErrTitleTooLong
	}

	if productType == "" {
		return nil, ErrTypeRequired
	}

	if !status.IsValid() {
		return nil, ErrInvalidStatus
	}

	if stockQuantity != nil && *stockQuantity < 0 {
		return nil, ErrNegativeStock
	}

	if weight != nil && *weight < 0 {
		return nil, ErrNegativeWeight
	}

	return &Product{
		Price:         price,
		Title:         title,
		Type:          productType,
		ImageURL:      imageURL,
		Status:        status,
		Color:         color,
		StockQuantity: stockQuantity,
		Description:   description,
		Weight:        weight,
		Rating:        rating,
		Size:          size,
	}, nil
}
