package user

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
)

type UserHandler struct {
	UserService *application.UserService
}

func (s *UserHandler) GetByPhone() http.HandlerFunc {
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
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).Warn("New user", "email:", user.Email)
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
