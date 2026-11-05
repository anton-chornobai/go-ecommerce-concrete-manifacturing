package product_handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/product/dto"
	"github.com/google/uuid"
)

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strID := r.PathValue("id")
	if strID == "" {
		h.logger.WarnContext(ctx, "Не вказано id ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		)
		http.Error(w, "Не вказано id ресурсу", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(strID)
	if err != nil {
		h.logger.WarnContext(ctx, "Невалідний UUID ресурсу",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.String("отримано", strID),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "id має бути валідним UUID", http.StatusBadRequest)
		return
	}

	h.logger.InfoContext(ctx, "Запит на оновлення продукту",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		slog.String("id", id.String()),
	)

	r.Body = http.MaxBytesReader(w, r.Body, maxBytesBodyLimit)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.logger.WarnContext(ctx, "Тіло запиту занадто велике при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.String("id", id.String()),
				slog.Int64("ліміт_байт", maxBytesErr.Limit),
			)
			http.Error(w, fmt.Sprintf("Тіло запиту занадто велике (максимум %dmb)", maxBytesBodyLimit>>20), http.StatusRequestEntityTooLarge)
		} else {
			h.logger.WarnContext(ctx, "Неправильний формат тіла запиту при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.String("id", id.String()),
				slog.String("помилка", err.Error()),
			)
			http.Error(w, "Неправильний формат запиту", http.StatusBadRequest)
		}
		return
	}

	parseForm := new(dto.ProductPatchRequest)
	if err := decoder.Decode(parseForm, r.PostForm); err != nil {
		h.logger.WarnContext(ctx, "Помилка декодування форми оновлення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.String("id", id.String()),
			slog.String("помилка", err.Error()),
		)
		http.Error(w, "Помилка декодування: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			file = nil
			header = nil
		case errors.Is(err, multipart.ErrMessageTooLarge):
			h.logger.WarnContext(ctx, "Файл завеликий при оновленні",
				slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
				slog.String("id", id.String()),
			)
			http.Error(w, "Файл завеликий", http.StatusRequestEntityTooLarge)
			return
		default:
			h.logger.WarnContext(ctx, "Помилка отримання файлу при оновленні",
				slog.String("id", id.String()),
				slog.String("помилка", err.Error()),
			)
			http.Error(w, "Неправильний запит", http.StatusBadRequest)
			return
		}
	}
	if file != nil {
		defer file.Close()
	}

	if err := h.ProductService.Update(ctx, id, *parseForm, file, header); err != nil {
		h.handleServiceError(w, r, err)
		return
	}

	h.logger.InfoContext(ctx, "Продукт успішно оновлено",
		slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
		slog.String("id", id.String()),
	)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Продукт з id %s успішно оновлено", id),
	}); err != nil {
		h.logger.ErrorContext(ctx, "Не вдалося закодувати відповідь після оновлення",
			slog.Group("source", slog.String("layer", "handler"), slog.String("func", "Update")),
			slog.String("id", id.String()),
			slog.String("помилка", err.Error()),
		)
	}
}
