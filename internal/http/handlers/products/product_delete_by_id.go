package product_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *ProductHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(strID)
	if err != nil {
		h.logger.WarnContext(ctx, "Невалідний UUID ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
			slog.String("отримано", strID),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "id має бути валідним UUID", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на видалення продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		slog.String("id", id.String()),
	)

	if err := h.ProductService.DeleteByID(ctx, id); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно видалено",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
		slog.String("id", id.String()),
	)

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Продукт з id %s успішно видалено", id),
	}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь після видалення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "DeleteByID")),
			slog.String("id", id.String()),
			slog.String("помилка", err.Error()),
		)
	}
}