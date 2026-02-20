package auth_handler

import (
	"encoding/json"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"net/http"
)

type AuthHandler struct {
	UserService *application.UserService
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password"`
}

func (s *AuthHandler) SignUpByEmail(w http.ResponseWriter, r *http.Request)  {
	var req SignupRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := s.UserService.Signup(req.Email, req.Password)

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
