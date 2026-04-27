package domain

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name          string
		price         int
		title         string
		productType   string
		status        ProductStatus
		stockQuantity *int
		weight        *int
		wantErr       bool
		errType       error
	}{
		{
			name:        "Price is 0",
			price:       0,
			title:       "Victory",
			productType: "tile",
			status:      ProductArchived,
			wantErr:     true,
			errType:     ErrInvalidPrice,
		},
		{
			name:        "Title length too short",
			price:       130,
			title:       "1",
			productType: "tile",
			status:      ProductArchived,
			wantErr:     true,
			errType:     ErrTitleTooShort,
		},
		{
			name:        "Title length too long",
			price:       130,
			title:       "kopdasokpasdosdakopdpakskodaspkoadspkdsapkdaspkpakdoskodsakdaspkaspkoadspkoadskpakdospkaskoadspko",
			productType: "tile",
			status:      ProductArchived,
			wantErr:     true,
			errType:     ErrTitleTooLong,
		},
		{
			name:        "productType is empty",
			price:       130,
			title:       "Plytaka",
			productType: "",
			status:      ProductArchived,
			wantErr:     true,
			errType:     ErrTypeRequired,
		},
		{
			name:        "Wrong product status",
			price:       200,
			title:       "Victory",
			productType: "tile",
			status:      ProductStatus("garbage_value"),
			wantErr:     true,
			errType:     ErrInvalidStatus,
		},
		{
			name:          "Wrong stock quantity",
			price:         200,
			title:         "Victory",
			productType:   "tile",
			status:        ProductArchived,
			stockQuantity: new(-10),
			wantErr:       true,
			errType:       ErrNegativeStock,
		},
		{
			name:          "Weight below zero",
			price:         200,
			title:         "Victory",
			productType:   "tile",
			status:        ProductArchived,
			weight: 		new(-19),
			wantErr:       true,
			errType:       ErrNegativeWeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewProduct(tt.price, tt.title, tt.productType, nil, tt.status, nil, tt.stockQuantity, nil, tt.weight, nil, nil)

			if tt.wantErr {
				if !errors.Is(err, tt.errType) {
					t.Errorf("got error %v, want %v", err, tt.errType)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
