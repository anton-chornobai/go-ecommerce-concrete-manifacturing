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

func (o *OrderService) Orders(limit int) {

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
