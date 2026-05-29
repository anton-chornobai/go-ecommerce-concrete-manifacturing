package middleware

import (
	"context"
	"errors"

	"net/http"

	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func GetUsersID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(w, "cookie not found", http.StatusUnauthorized)
				return
			}
			http.Error(w, "error reading cookie", http.StatusBadRequest)
			return
		}

		token := cookie.Value

		id, err := jwtmanager.GetUsersID(token)

		if err != nil {
			http.Error(w, "problem with token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, id)

		newRequestContext := r.WithContext(ctx)

		next.ServeHTTP(w, newRequestContext)
	}
}
