package application

import "github.com/anton-chornobai/beton.git/internal/modules/orders/domain"

type OrderService struct {
	repo domain.OrderRepository
}

type CreateOrderRequest struct {
	UserID   string
	Items    []domain.OrderItem
	Discount int
}

func (o *OrderService) Orders(limit int) {

}

func (o *OrderService) Create(req CreateOrderRequest) (*domain.Order, error) {
	order, err := domain.NewOrder(req.UserID, req.Items, req.Discount)

	if err != nil {
		return nil, nil
	}

	err = o.repo.Save(order); 

	if err != nil {
		return nil, err;
	}

	return order, nil
}
