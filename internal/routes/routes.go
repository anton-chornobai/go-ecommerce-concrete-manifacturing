package routes

import (
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/handlers"
	"github.com/anton-chornobai/beton.git/internal/middleware"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
)


func SetUpRouter(userAppService *application.UserAppService, ordersAppService *ordersApp.OrderService) http.Handler  {
	usersHandler := handlers.UsersHandler {
		UserService: userAppService,
	}

	ordersHandler := handlers.OrdersHandler {
		OrdersService: ordersAppService,
	}
	router := http.NewServeMux()

	router.HandleFunc("POST /auth", usersHandler.Register())
	router.Handle("GET /profile", middleware.VerifyToken(handlers.GetProfile()))
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Products"))
	})
	router.HandleFunc("GET /user", usersHandler.GetByPhone())
    router.HandleFunc("GET /admin", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Admin root"))
    })
	
	router.HandleFunc("POST /orders", ordersHandler.Create())

	return  middleware.CorsMiddleware(router)
}