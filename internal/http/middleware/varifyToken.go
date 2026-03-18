package middleware

import (
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, "unautherized", http.StatusUnauthorized)
			return
		}

		claims, err := jwtmanager.ValidateToken(cookie.Value)

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := jwtmanager.AddClaimsToContext(r.Context(), claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(userService *application.UserService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, "cookie" + err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := jwtmanager.ValidateToken(cookie.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		id, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "could not get id", http.StatusForbidden)
			return
		}

		isAdmin, err := userService.IsAdmin(id)

		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
