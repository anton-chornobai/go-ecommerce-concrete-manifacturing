package routes

import (
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/handlers"
	"github.com/anton-chornobai/beton.git/internal/middleware"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
)


func SetUpRouter(userServives *application.UserAppService) http.Handler  {
	ordersHadnler := &handlers.OrdersHandler{}
	router := http.NewServeMux()

	router.HandleFunc("POST /auth", handlers.RegisterUser(userServives))
	router.Handle("GET /profile", middleware.VerifyToken(handlers.GetProfile()))
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Products"))
	})

    router.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Admin root"))
    })
	

	router.HandleFunc("GET /admin/orders", ordersHadnler.GetOrders())

	

	return  middleware.CorsMiddleware(router)
}