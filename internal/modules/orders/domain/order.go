package domain

import (
	"errors"
	"time"
)

type Size struct {
	Height    int
	Width     int
	Thickness int
}

type OrderItem struct {
	ProductID int    `db:"product_id" json:"product_id"`
	Title     string `db:"title" json:"title"`
	UnitPrice int    `db:"unit_price" json:"unit_price"`
	Type      string `db:"type" json:"type"`
	Quantity  int    `db:"quantity" json:"quantity"`
	Color     string `db:"color" json:"color"`
	Height    int    `db:"height" json:"height"`
	Width     int    `db:"width" json:"width"`
	Material  string `db:"material" json:"material"`
	Thickness int    `db:"thickness" json:"thickness"`
}

type Order struct {
	ID        int         `json:"id"`
	UserID    string      `json:"user_id"`
	Name      string      `json:"name"`
	Items     []OrderItem `json:"items"`
	Total     int         `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	Discount  int         `json:"discount"`
}


func NewOrderService()  {

}


func NewOrder(userId string, items []OrderItem, discount int) (*Order, error) {
	if userId == "" {
		return nil, errors.New("userID is required")
	}

	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	total := 0
	for _, item := range items {
		total += item.UnitPrice * item.Quantity
	}

	total -= total * discount

	return &Order{
		UserID:     userId,
		Items:      items,
		Total: total,
		Status:     "created",
		CreatedAt:  time.Now(),
	}, nil
}
