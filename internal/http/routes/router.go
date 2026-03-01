package routes

import (
	"net/http"

	// "github.com/anton-chornobai/beton.git/internal/http/handlers"
	"github.com/anton-chornobai/beton.git/internal/http/handlers"
	user_handler "github.com/anton-chornobai/beton.git/internal/http/handlers/user"

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
	authHandler := user_handler.AuthHandler{
		UserService: userService,
	}

	productHandler := handlers.ProductHandler{
		ProductService: productService,
	}

	userHandler := user_handler.UserHandler{
		UserService: userService,
	}

	// orderHandler := handlers.OrdersHandler {
	// 	OrdersService: orderService,
	// }

	router := http.NewServeMux()
	router.HandleFunc("POST /signup", authHandler.SignupByEmail)
	router.HandleFunc("POST /login", authHandler.LoginByEmail)
	router.HandleFunc("POST /verify", authHandler.Verify)
	router.HandleFunc("GET /profile", userHandler.GetByID)

	router.Handle(
		"POST /admin/products",
		middleware.AdminOnly(http.HandlerFunc(productHandler.Add)),
	)
	// router.HandleFunc("POST /admin/products", )

	return middleware.CorsMiddleware(router)
}
