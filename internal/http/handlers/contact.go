package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/anton-chornobai/beton.git/internal/modules/contact/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/dto"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/service"
)

type UserContactHandler struct {
	UserContactService *service.UserContactService
	logger             *slog.Logger
}

func NewUserContactHandler(service *service.UserContactService, logger *slog.Logger) *UserContactHandler {
	return &UserContactHandler{
		UserContactService: service,
		logger:             logger,
	}
}

func (h *UserContactHandler) Post(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var req dto.UserContactPostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("UserContactHandler.Post", "err", err.Error())
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Email == "" || req.Message == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	err := h.UserContactService.Post(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmailTooLong),
			errors.Is(err, domain.ErrInvalidNumber),
			errors.Is(err, domain.ErrInvalidNumberSymbol),
			errors.Is(err, domain.ErrWrongEmailFormat),
			errors.Is(err, domain.ErrNameTooLong):

			http.Error(w, err.Error(), http.StatusBadRequest)
			h.logger.Warn("UserContactHandler.Post", "err", err.Error())

			return
		default:
			h.logger.Error("UserContactHandler.Post", "err", err.Error())
			http.Error(w, "Щось пішло не так...", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Новий контакт створено!",
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *UserContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	id := r.PathValue("id")
	if id == "" {
		h.logger.Warn("UserContactHandler.Delete", "err", "missing id")
		http.Error(w, "id not provided", http.StatusBadRequest)
		return
	}

	if id == "" {
		h.logger.Warn("UserContactHandler.Delete", "err", "id not provided")
		http.Error(w, "id not provided", http.StatusBadRequest)
		return
	}

	err := h.UserContactService.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			h.logger.Warn("UserContactHandler.Delete", "err", err.Error())
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		h.logger.Warn("UserContactHandler.Delete", "err", err.Error())
			http.Error(w, "Щось пішло не так", http.StatusInternalServerError)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Видалено",
	})
}
