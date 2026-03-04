package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersHandler struct {
	OrdersService *application.OrderService
}

func (o *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var order *domain.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = o.OrdersService.Create(ctx, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}
