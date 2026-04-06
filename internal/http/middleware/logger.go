package middleware

import (
	"log/slog"
	"net/http"
)

func LogMethodInfo(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("METHOD", "endpoint", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
