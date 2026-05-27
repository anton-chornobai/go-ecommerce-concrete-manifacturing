package domain

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func TestNewProduct(t *testing.T) {

	tests := []struct {
		name            string
		price           int
		title           string
		productType     string
		status          ProductStatus
		stockQuantity   *int
		weight          *int
		size            *Size
		wantErr         bool
		expectedErrType error
	}{
		{
			name:            "Price is 0",
			price:           0,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         true,
			expectedErrType: ErrInvalidPrice,
		},
		{
			name:            "Price is negative",
			price:           -100,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         true,
			expectedErrType: ErrInvalidPrice,
		},
		{
			name:            "Valid price",
			price:           100,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         false,
			expectedErrType: nil,
		},
		{
			name:            "Title length too short",
			price:           130,
			title:           "1",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         true,
			expectedErrType: ErrTitleTooShort,
		},
		{
			name:        "Valid random title within length bounds",
			price:       130,
			title:       generateRandomTitle(),
			productType: "tile",
			status:      ProductArchived,
			wantErr:     false,
		},
		{
			name:            "Title length too long",
			price:           130,
			title:           "This title is definitely over thirty chars long",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         true,
			expectedErrType: ErrTitleTooLong,
		},
		{
			name:            "Valid title length",
			price:           130,
			title:           "Valid Title",
			productType:     "tile",
			status:          ProductArchived,
			wantErr:         false,
			expectedErrType: nil,
		},
		{
			name:            "productType is empty",
			price:           130,
			title:           "Plytaka",
			productType:     "",
			status:          ProductArchived,
			wantErr:         true,
			expectedErrType: ErrTypeRequired,
		},
		{
			name:            "Wrong product status",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductStatus("garbage_value"),
			wantErr:         true,
			expectedErrType: ErrInvalidStatus,
		},
		{
			name:            "Wrong stock quantity",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			stockQuantity:   new(-10),
			wantErr:         true,
			expectedErrType: ErrNegativeStock,
		},
		{
			name:            "Weight below zero",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			weight:          new(-19),
			wantErr:         true,
			expectedErrType: ErrNegativeWeight,
		},
		{
			name:            "Negative width is not valid",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			weight:          new(10),
			size:            &Size{Width: -10, Height: 10},
			wantErr:         true,
			expectedErrType: ErrNegativeWidth,
		},
		{
			name:            "Negative height is not valid",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			weight:          new(10),
			size:            &Size{Width: 10, Height: -10},
			wantErr:         true,
			expectedErrType: ErrNegativeHeight,
		},
		{
			name:            "Valid size width and height",
			price:           200,
			title:           "Victory",
			productType:     "tile",
			status:          ProductArchived,
			weight:          new(10),
			size:            &Size{Width: 10, Height: 10},
			wantErr:         false,
			expectedErrType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProduct(tt.price, tt.title, tt.productType, nil, tt.status, nil, tt.stockQuantity, nil, tt.weight, nil, tt.size)

			if tt.wantErr {
				if !errors.Is(err, tt.expectedErrType) {
					t.Errorf("got error %v, want %v", err, tt.expectedErrType)
				}
				if got != nil {
					t.Errorf("expected nil product on error, got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if got == nil {
					t.Fatal("expected product to be created, but got nil")
				}
				if got.Price != tt.price {
					t.Errorf("Field Price mapping mismatch: got %v, want %v", got.Price, tt.price)
				}
				if got.Title != tt.title {
					t.Errorf("Field Title mapping mismatch: got %v, want %v", got.Title, tt.title)
				}
				if got.Type != tt.productType {
					t.Errorf("Field Type mapping mismatch: got %v, want %v", got.Type, tt.productType)
				}
				if got.Status != tt.status {
					t.Errorf("Field Status mapping mismatch: got %v, want %v", got.Status, tt.status)
				}

				if tt.stockQuantity != nil && (got.StockQuantity == nil || *got.StockQuantity != *tt.stockQuantity) {
					t.Errorf("Field StockQuantity mismatch")
				}
				if tt.weight != nil && (got.Weight == nil || *got.Weight != *tt.weight) {
					t.Errorf("Field Weight mismatch")
				}
				if tt.size != nil && (got.Size == nil || got.Size.Width != tt.size.Width || got.Size.Height != tt.size.Height) {
					t.Errorf("Field Size mismatch")
				}
			}
		})
	}
}

func generateRandomTitle() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	length := r.Intn(MaxTitleLength-MinTitleLength+1) + MinTitleLength

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
