package models


type Size struct {
    Width  float64 `json:"width"`
    Height float64 `json:"height"`
}

type Product struct {
    ID            int      `json:"id"`
    Title         string   `json:"title"`
    Price         float64  `json:"price"`
    Type          string   `json:"type"`
    ImageURL      string   `json:"imageUrl"`
    Color         string   `json:"color"`
    StockQuantity *int     `json:"stockQuantity,omitempty"`
    Description   *string  `json:"description,omitempty"`
    Weight        *float64 `json:"weight,omitempty"`
    Rating        *float64 `json:"rating,omitempty"`
    Size          *Size    `json:"size,omitempty"`
}