package auth_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
)

type AuthHandler struct {
	UserService *application.UserService
}

type SignupEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type LoginEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

// Yet to be implemented
type SignupNumberRequest struct {
	Number string `json:"number"`
}

type LoginNumberRequest struct {
	Number string `json:"number"`
}

func (s *AuthHandler) SignupByEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req SignupEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := s.UserService.SignupByEmail(ctx, req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // set to true on production!!!
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/signup",
	})

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	}); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *AuthHandler) LoginByEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second * 5)
	defer cancel();

	var req LoginEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload credentials", http.StatusBadRequest)
	}

	token, err := s.UserService.LoginByEmail(ctx, req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // set to true on production!!!
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/login",
	})

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)


	if err := json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	}); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
} 
