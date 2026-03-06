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
	orders, err := o.Repo.Get(ctx, limit);

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderService) Create(ctx context.Context, req *domain.Order) error {
	order, err := domain.NewOrder(req.UserID, req.OrderName, req.Items, req.Discount, req.ShippingAddress, req.ShippingCity, req.ShippingPostalCode)

	if err != nil {
		return err
	}
	err = o.Repo.Create(ctx, order)

	if err != nil {
		return err
	}

	return nil
}
