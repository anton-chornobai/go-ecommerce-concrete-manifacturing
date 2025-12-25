package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
)


func RegisterUser(userService *application.UserAppService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var number domain.AuthenticationUserRequest

		if err := json.NewDecoder(r.Body).Decode(&number); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
		}

		registerResult, err := userService.Register(number) 

		if err != nil {
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