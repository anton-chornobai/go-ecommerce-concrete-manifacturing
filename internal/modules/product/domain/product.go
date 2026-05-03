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
	ErrProductNotFound = errors.New("продукт не знайдено")
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

type ProductPatch struct {
    Price         *int
    Title         *string
    Type          *string
    Status        *ProductStatus
    Color         *string
    Description   *string
    StockQuantity *int
    Weight        *int
    Rating        *int
    SizeWidth     *int
    SizeHeight    *int
    ImageURL      *string
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
	p := &Product{
		Price:         price,
		Title:         title,
		Type:          productType,
		Color:         color,
		Status:        status,
		ImageURL:      imageURL,
		StockQuantity: stockQuantity,
		Description:   description,
		Weight:        weight,
		Rating:        rating,
		Size:          size,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Product) Validate() error {
	if p.Price <= 0 {
		return ErrInvalidPrice
	}

	titleLen := utf8.RuneCountInString(p.Title)
	if titleLen < MinTitleLength {
		return ErrTitleTooShort
	}
	if titleLen > MaxTitleLength {
		return ErrTitleTooLong
	}

	if p.Type == "" {
		return ErrTypeRequired
	}

	if !p.Status.IsValid() {
		return ErrInvalidStatus
	}

	if p.StockQuantity != nil && *p.StockQuantity < 0 {
		return ErrNegativeStock
	}

	if p.Weight != nil && *p.Weight < 0 {
		return ErrNegativeWeight
	}

	return nil
}

func (p *ProductPatch) Validate() error {
    if p.Price != nil && *p.Price <= 0 {
        return ErrInvalidPrice
    }

    if p.Title != nil {
        titleLen := utf8.RuneCountInString(*p.Title)
        if titleLen < MinTitleLength {
            return ErrTitleTooShort
        }
        if titleLen > MaxTitleLength {
            return ErrTitleTooLong
        }
    }

    if p.Type != nil && *p.Type == "" {
        return ErrTypeRequired
    }

    if p.Status != nil && !p.Status.IsValid() {
        return ErrInvalidStatus
    }

    if p.StockQuantity != nil && *p.StockQuantity < 0 {
        return ErrNegativeStock
    }

    if p.Weight != nil && *p.Weight < 0 {
        return ErrNegativeWeight
    }

    return nil
}