package application

import (
	"log/slog"

	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrderService struct {
	Repo domain.OrderRepository
	log *slog.Logger
}

type CreateOrderRequest struct {
	UserID   string             `json:"user_id"`
	Name     string             `json:"name"`
	Items    []domain.OrderItem `json:"items"`
	Discount int                `json:"discount"`
}

func NewOrderService(repo domain.OrderRepository, logger *slog.Logger) *OrderService {
	return &OrderService{
		Repo: repo,
		log: logger,
	}
}

func (o *OrderService) Orders(limit int) {

}

func (o *OrderService) Create(req CreateOrderRequest) (*domain.Order, error) {
	const op = "OrderService.Create";

	log := o.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to create order", "user_id", req.UserID)

	order, err := domain.NewOrder(req.UserID, req.Items, req.Discount)

	if err != nil {
		log.Warn("Couldnt create order")
		return nil, err
	}

	err = o.Repo.Save(order)

	if err != nil {
		log.Warn("Couldnt save order to db", "user_id", req.UserID)
		return nil, err
	}

	log.Warn("Order created", "user_id", req.UserID, "name", req.Name)

	return order, nil
}
