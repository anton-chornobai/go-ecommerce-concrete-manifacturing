package user

import (
	"context"
	"encoding/json"
	"errors"

	"net/http"
	"time"

	// "github.com/anton-chornobai/beton.git/internal/mail"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
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

	// Decode request body
	var req SignupEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	err := s.UserService.SignupByEmail(ctx, req.Email, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		if errors.Is(err, infra.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "This user already exists!",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "signup failed: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "signup successful, please check your email to verify",
	})
}
func (s *AuthHandler) LoginByEmail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var req LoginEmailRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid payload credentials", http.StatusBadRequest)
	}

	token, err := s.UserService.LoginByEmail(ctx, req.Email, req.Password)

	if err != nil {
		if errors.Is(err, application.ErrAccountNotVerified) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden) // 403
			json.NewEncoder(w).Encode(map[string]string{
				"error": "account not verified",
			})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // set to true on production!!!
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
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

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.UserService.VerifyUser(ctx, req.Email, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // set true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "email verified successfully",
		"token":   token,
	})
}
