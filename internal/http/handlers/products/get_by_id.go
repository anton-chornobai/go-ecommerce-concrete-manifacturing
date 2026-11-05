package product_handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(strID)
	if err != nil {
		h.logger.WarnContext(ctx, "Невалідний UUID ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
			slog.String("отримано", strID),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "id має бути валідним UUID", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на отримання продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		slog.String("id", id.String()),
	)

	product, err := h.ProductService.GetById(ctx, id)
	if err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно отримано",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
		slog.String("id", id.String()),
	)

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(map[string]any{"data": product}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProductByID")),
			slog.String("id", id.String()),
			slog.String("помилка", err.Error()),
		)
	}
}