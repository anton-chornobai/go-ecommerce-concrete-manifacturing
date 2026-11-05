package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	MaxTitleLength = 30
	MinTitleLength = 2
)

var (
	ErrInvalidPrice       = errors.New("ціна повинна бути більшою за 0")
	ErrTitleTooShort      = errors.New("назва повинна містити принаймні 2 символи")
	ErrTitleTooLong       = errors.New("назва занадто довга (максимум 30 символів)")
	ErrTypeRequired       = errors.New("тип продукту є обов'язковим")
	ErrInvalidStatus      = errors.New("недопустимий статус продукту")
	ErrNegativeStock      = errors.New("кількість на складі не може бути від'ємною")
	ErrNegativeWeight     = errors.New("вага не може бути від'ємною")
	ErrProductNotFound    = errors.New("продукт не знайдено")
	ErrTitleAlreadyExists = errors.New("назва продукту вже існує")
	ErrNegativeWidth      = errors.New("ширина не може бути від'ємною")
	ErrNegativeHeight     = errors.New("висота не може бути від'ємною")

	ErrNegativeImagePosition = errors.New("позиція не може бути негативною")
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

type Product struct {
	ID            uuid.UUID
	Price         int
	Title         string
	Type          string
	Color         *string
	Status        ProductStatus
	ImageURLs     []ProductImage
	Description   *string
	StockQuantity *int
	Weight        *int
	Rating        *int
	Size          *Size
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ProductImage struct {
	ID        uuid.UUID
	Position  int
	Key       string
	CreatedAt time.Time
}
type Size struct {
	Width  int
	Height int
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

func NewProductImage(position int, url string) (*ProductImage, error) {
	if position < 0 {
		return nil, ErrNegativeImagePosition
	}

	return &ProductImage{
		uuid.New(),
		position,
		url,
		time.Now(),
	}, nil
}

func NewProduct(
	price int,
	title string,
	productType string,
	color *string,
	status ProductStatus,
	imageURLs []ProductImage,
	stockQuantity *int,
	description *string,
	weight *int,
	rating *int,
	size *Size,
) (*Product, error) {
	p := &Product{
		ID:            uuid.New(),
		Price:         price,
		Title:         title,
		Type:          productType,
		Color:         color,
		Status:        status,
		ImageURLs:     imageURLs,
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

	if p.Size != nil {
		if p.Size.Width < 0 {
			return ErrNegativeWidth
		}
		if p.Size.Height < 0 {
			return ErrNegativeHeight
		}
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

	if p.SizeWidth != nil && *p.SizeWidth < 0 {
		return ErrNegativeWidth
	}
	if p.SizeHeight != nil && *p.SizeHeight < 0 {
		return ErrNegativeHeight
	}

	return nil
}
