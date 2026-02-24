package routes

import (
	"net/http"

	// "github.com/anton-chornobai/beton.git/internal/http/handlers"
	auth_handler "github.com/anton-chornobai/beton.git/internal/http/handlers/auth"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	userService "github.com/anton-chornobai/beton.git/internal/modules/user/application"
)


func SetUpRoutes(userService *userService.UserService, orderService *application.OrderService) *http.ServeMux { 
	authHandler := auth_handler.AuthHandler{
		UserService: userService,
	}

	// userHandler := handlers.UserHandler{
	// 	UserService: userService,
	// }

	// orderHandler := handlers.OrdersHandler {
	// 	OrdersService: orderService,
	// }


	router := http.NewServeMux()
	router.HandleFunc("POST /auth/signup", authHandler.SignupByEmail)
	router.HandleFunc("POST /auth/login", authHandler.LoginByEmail)
	// router.HandleFunc("GET /admin/products", middleware.AdminOnly())
	// router.HandleFunc("POST /admin/products", )

	return  router;
}