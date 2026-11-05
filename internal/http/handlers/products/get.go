package product_handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	userMiddleware "github.com/anton-chornobai/beton.git/internal/http/middleware/user"
	"github.com/anton-chornobai/beton.git/internal/modules/product/domain"
)

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	h.logger.InfoContext(ctx, "Запит на отримання списку продуктів",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
	)

	role, _ := r.Context().Value(userMiddleware.RoleContextKey).(string)

	var status *domain.ProductStatus
	statusVal := r.URL.Query().Get("status")
	if statusVal == string(domain.ProductArchived) || statusVal == string(domain.ProductDisplayed) {
		s := domain.ProductStatus(statusVal)

		if s == domain.ProductArchived && role != "admin" {
			h.logger.WarnContext(ctx, "Спроба доступу до архівованих продуктів без прав",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
				slog.String("role", role),
			)
			http.Error(w, "Доступ заборонено", http.StatusForbidden)
			return
		}

		status = &s
	}

	products, err := h.ProductService.GetProducts(ctx, 20, status)
	if err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося отримати список продуктів",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Щось пішло не так", http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(ctx, "Список продуктів успішно отримано",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
		slog.Int("кількість", len(products)),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string][]domain.Product{"data": products}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "GetProducts")),
			slog.String("помилка", err.Error()),
		)
	}
}