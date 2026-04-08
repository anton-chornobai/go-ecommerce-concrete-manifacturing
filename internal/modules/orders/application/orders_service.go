package application

import (
	"context"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
	"log/slog"
)

type OrderService struct {
	Repo domain.OrderRepository
	Log  *slog.Logger
}

func NewOrderService(repo domain.OrderRepository, log *slog.Logger) *OrderService {
	return &OrderService{
		Repo: repo,
		Log:  log,
	}
}

func (o *OrderService) Get(ctx context.Context, limit int) ([]domain.Order, error) {
	orders, err := o.Repo.Get(ctx, limit)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderService) Create(ctx context.Context, req *domain.Order) (int, error) {
	order, err := domain.NewOrder(req.UserID, req.OrderName, req.CustomerName, req.Items, req.Discount, req.CustomerNumber, req.ShippingAddress, req.ShippingCity, req.ShippingPostalCode)

	if err != nil {
		return 0, err
	}
	id, err := o.Repo.Create(ctx, order)

	if err != nil {
		return 0, err
	}

	return  id, nil
}
