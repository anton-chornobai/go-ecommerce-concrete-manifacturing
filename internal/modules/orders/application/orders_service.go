package application

import (
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrderService struct {
	Repo domain.OrderRepository
}

type CreateOrderRequest struct {
	UserID   string             `json:"user_id"`
	Name     string             `json:"name"`
	Items    []domain.OrderItem `json:"items"`
	Discount int                `json:"discount"`
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{
		Repo: repo,
	}
}

func (o *OrderService) Orders(limit int) {

}

func (o *OrderService) Create(req CreateOrderRequest) (*domain.Order, error) {
	order, err := domain.NewOrder(req.UserID, req.Items, req.Discount)

	if err != nil {
		return nil, err
	}

	err = o.Repo.Save(order)

	if err != nil {
		return nil, err
	}

	return order, nil
}
