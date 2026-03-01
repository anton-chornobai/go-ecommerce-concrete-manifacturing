package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
)

type UserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Email    string  `json:"email"`
	Role     string  `json:"role"`
	Address  string  `json:"address"`
	Number   *string `json:"number,omitempty"`
	Verified bool   `json:"is_verified"`
}

type UserHandler struct {
	UserService *application.UserService
}

func (s *UserHandler) GetByPhone() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var number struct{ Number string }

		err := json.NewDecoder(r.Body).Decode(&number)
		if err != nil {
			http.Error(w, "invalid argument", http.StatusBadRequest)
			return
		}

		user, err := s.UserService.GetByPhone(number.Number)
		if err != nil {
			log.Println(err)
			http.Error(w, "invalid argument", http.StatusInternalServerError)
			return
		}
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Warn("New user", "email:", user.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]any{
			"user":    user,
			"message": "User",
		}); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "No session cookie", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	jwt := cookie.Value

	user, err := s.UserService.GetByID(jwt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Println("GetByID error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     user.Role,
		Address:  user.Address,
		Number:   user.Number,
		Verified: user.IsVerified,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"user":    resp,
		"message": "User retrieved successfully",
	}); err != nil {
		log.Println("JSON encode error:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
