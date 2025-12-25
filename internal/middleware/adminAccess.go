package middleware

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/anton-chornobai/beton.git/internal/utils"
// )

// func VerifyAdminAccess(http.Handler) http.Handler  {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
// 		cookie, err := r.Cookie("jwt")

// 		if err != nil {
// 			http.Error(w, "Is not authothenticated", http.StatusUnauthorized)
// 		}

// 		claims, err := utils.GetUserClaims(cookie.Value, os.Getenv("SECRET"))

// 		fmt.Println(claims.Role)
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("ACCESSED ADMIN"))
// 	})
// }