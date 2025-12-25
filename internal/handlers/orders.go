package handlers

import "net/http"


type OrdersHandler struct {

}

func (r *OrdersHandler) GetOrders() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("You hit the orders"))
	})
}