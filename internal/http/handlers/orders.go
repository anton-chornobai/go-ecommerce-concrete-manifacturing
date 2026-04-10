package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anton-chornobai/beton.git/internal/http/middleware"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/domain"
)

type OrdersHandler struct {
	log           *slog.Logger
	OrdersService *application.OrderService
}

func NewOrdersHandler(log *slog.Logger, orderService *application.OrderService) *OrdersHandler {
	return &OrdersHandler{log: log, OrdersService: orderService}
}

func (o *OrdersHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	limitStr := r.URL.Query().Get("limit")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "No orders found",
			"data":    []any{},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(map[string]any{
		"message": "Orders",
		"data":    orders,
	})
	if err != nil {
		o.log.Info("encode error:", fmt.Sprintf("er %w", err))
	}
}

func (o *OrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	defer r.Body.Close()
	var order domain.Order

	issuerId, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok || issuerId == "" {
		fmt.Println(issuerId)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Setting ID of the user who Created the Order
	order.UserID = issuerId

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := o.OrdersService.Create(ctx, &order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"id": id,
	}); err != nil {
		o.log.Error("encode error", "err", err)
	}
}

func (o *OrdersHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	partsOfURL := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	strID := partsOfURL[len(partsOfURL)-1]
	intID, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Failed to convert string id", http.StatusBadRequest)
		return
	}

	err = o.OrdersService.Delete(ctx, intID)
	if err != nil {
		o.log.Warn("OrderService.Delete", "ERR:", err)
		http.Error(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "successfully deleted",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
