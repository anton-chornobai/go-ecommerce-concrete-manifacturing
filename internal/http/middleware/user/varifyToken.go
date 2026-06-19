package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

type contextKey string

const RoleContextKey contextKey = "role"

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
			http.Error(w, "cookie"+err.Error(), http.StatusUnauthorized)
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
			if errors.Is(err, domain.ErrUnauthorized) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Error(w, "щось пішло не так", http.StatusInternalServerError)
			return
		}

		if !isAdmin {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetUsersRoleWithContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			ctx := context.WithValue(r.Context(), RoleContextKey, "user")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, err := jwtmanager.ParseToken(cookie.Value)
		if err != nil {
			ctx := context.WithValue(r.Context(), RoleContextKey, "user")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role == "" {

			role = "user"
		}

		ctx := context.WithValue(r.Context(), RoleContextKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
