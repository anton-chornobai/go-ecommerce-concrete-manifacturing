package routes

import (
	"net/http"

	// "github.com/anton-chornobai/beton.git/internal/http/handlers"
	"github.com/anton-chornobai/beton.git/internal/http/handlers"
	auth_handler "github.com/anton-chornobai/beton.git/internal/http/handlers/user"
	"github.com/anton-chornobai/beton.git/internal/http/middleware"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	productService "github.com/anton-chornobai/beton.git/internal/modules/product/application"
	userService "github.com/anton-chornobai/beton.git/internal/modules/user/application"
)

func SetUpRoutes(
	userService *userService.UserService,
	orderService *application.OrderService,
	productService productService.ProductService,
) http.Handler {
	authHandler := auth_handler.AuthHandler{
		UserService: userService,
	}

	productHandler := handlers.ProductHandler{
		ProductService: productService,
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
	router.Handle(
		"POST /admin/products",
		middleware.AdminOnly(http.HandlerFunc(productHandler.Add)),
	)
	// router.HandleFunc("POST /admin/products", )

	return middleware.CorsMiddleware(router)
}
