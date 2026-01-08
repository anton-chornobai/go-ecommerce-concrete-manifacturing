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
	ProductID int    `db:"product_id"`
	Title     string `db:"title"`
	UnitPrice int    `db:"unit_price"`
	Type      string `db:"type"`
	Quantity  int    `db:"quantity"`
	Color     string `db:"color"`
	Height    int    `db:"height"`
	Width     int    `db:"width"`
	Material  string `db:"material"`
	Thickness int    `db:"thickness"`
}
type Order struct {
	ID        int
	UserID    string
	Name 	  string
	Items     []OrderItem
	Total     int
	Status    string
	CreatedAt time.Time
	Discount  int
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
