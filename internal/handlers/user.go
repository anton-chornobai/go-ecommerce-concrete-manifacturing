package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)

type UsersHandler struct {
	UserService *application.UserService
}

func (s *UsersHandler) Register() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var number domain.AuthenticationUserRequest

		if err := json.NewDecoder(r.Body).Decode(&number); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		registerResult, err := s.UserService.Register(number)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    registerResult.Token,
			HttpOnly: true,
			Secure:   false, // set to true on production!!!
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]any{
			"user":    registerResult.User,
			"message": "User created successfully",
		}); err != nil {
			http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (s *UsersHandler) GetByPhone() http.HandlerFunc {
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
