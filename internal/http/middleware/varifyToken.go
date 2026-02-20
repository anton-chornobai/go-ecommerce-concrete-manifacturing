package middleware

import (
	"net/http"

	jwtmanager "github.com/anton-chornobai/beton.git/internal/lib/jwt"
)

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		if err != nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
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
