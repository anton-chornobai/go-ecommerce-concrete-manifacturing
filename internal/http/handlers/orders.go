package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersHandler struct {
	OrdersService *application.OrderService
}

func (o *OrdersHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	
	limitStr := r.URL.Query().Get("limit")
	//default limit
	limit := 10 
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = l
	}

	orders, err := o.OrdersService.Get(ctx, limit)
	if err != nil {
		http.Error(w, "Internal", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK);
	w.Header().Set("Content-Type", "application/json")  
	if err := json.NewEncoder(w).Encode(map[string]any { 
		"message": "Order successfuly created!",
		"data": orders,
	}); err != nil {
		http.Error(w, "Failed to send data", http.StatusInternalServerError)
	} 
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

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(200)
}
