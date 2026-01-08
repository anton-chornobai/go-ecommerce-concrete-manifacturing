package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
)

type OrdersHandler struct {
	OrdersService *application.OrderService
}

func (o *OrdersHandler) Create() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var order application.CreateOrderRequest

		err := json.NewDecoder(r.Body).Decode(&order)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdOrder, err := o.OrdersService.Create(order)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		if err := json.NewEncoder(w).Encode(createdOrder); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
