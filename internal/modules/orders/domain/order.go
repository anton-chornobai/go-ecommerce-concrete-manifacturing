package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type PaymentStatus string

const (
	OrderPending   OrderStatus = "очікується"
	OrderConfirmed OrderStatus = "підтверджено"
	OrderShipped   OrderStatus = "відправлено"
	OrderDelivered OrderStatus = "доставлено"
	OrderCancelled OrderStatus = "скасовано"
)

const (
	PaymentUnpaid   PaymentStatus = "неоплачено"
	PaymentPaid     PaymentStatus = "оплачено"
	PaymentFailed   PaymentStatus = "не вдалося"
	PaymentRefunded PaymentStatus = "повернено"
)

type Order struct {
	ID                 int           `json:"id"`
	Total              int           `json:"total"`
	CustomerName       string        `json:"customer_name"`
	UserID             string        `json:"user_id"`
	OrderName          string        `json:"order_name"`
	Items              []OrderItem   `json:"items"`
	Status             OrderStatus   `json:"status"`
	PaymentStatus      PaymentStatus `json:"payment_status"`
	CustomerNumber     *string       `json:"customer_number,omitempty"`
	Discount           *int          `json:"discount,omitempty"`
	ShippingAddress    *string       `json:"shipping_address,omitempty"`
	ShippingCity       *string       `json:"shipping_city,omitempty"`
	ShippingPostalCode *string       `json:"shipping_postal_code,omitempty"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

type Size struct {
	Height    int `json:"height"`
	Width     int `json:"width"`
	Thickness int `json:"thickness"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	Material  string    `json:"material"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice int       `json:"unit_price"`
	Size      Size      `json:"size"`
}

func NewOrder(
	userID string,
	orderName string,
	customerName string,
	items []OrderItem,
	discount *int,
	number *string,
	shippingAddress *string,
	shippingCity *string,
	shippingPostalCode *string,
) (*Order, error) {
	if customerName == "" {
		return nil, errors.New("customer name required")
	}
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	if orderName == "" {
		return nil, errors.New("order name is required")
	}

	if len(items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	if number == nil || *number == "" || len(*number) < 7 || len(*number) > 12 {
		return nil, errors.New("invalid number")
	}

	if shippingAddress == nil || *shippingAddress == "" {
		return nil, errors.New("shipping address is required")
	}
	if shippingCity == nil || *shippingCity == "" {
		return nil, errors.New("shipping city is required")
	}
	if shippingPostalCode == nil || *shippingPostalCode == "" {
		return nil, errors.New("shipping postal code is required")
	}
	if discount != nil {
		if *discount < 0 || *discount > 100 {
			return nil, errors.New("discount must be between 0 and 100")
		}
	}

	total := 0

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for item %s", item.Title)
		}
		if item.UnitPrice < 0 {
			return nil, fmt.Errorf("invalid unit price for item %s", item.Title)
		}

		total += item.UnitPrice * item.Quantity
	}

	if discount != nil && *discount > 0 {
		discountAmount := total * (*discount) / 100
		total -= discountAmount
	}

	now := time.Now()

	order := &Order{
		UserID:             userID,
		CustomerName:       customerName,
		OrderName:          orderName,
		Items:              items,
		Total:              total,
		Status:             OrderPending,
		PaymentStatus:      PaymentUnpaid,
		Discount:           discount,
		CustomerNumber:     number,
		ShippingAddress:    shippingAddress,
		ShippingCity:       shippingCity,
		ShippingPostalCode: shippingPostalCode,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return order, nil
}
