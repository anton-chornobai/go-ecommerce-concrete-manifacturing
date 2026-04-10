package routes

import (
	"log/slog"
	"net/http"

	// "github.com/anton-chornobai/beton.git/internal/http/handlers"
	"github.com/anton-chornobai/beton.git/internal/http/handlers"
	user_handler "github.com/anton-chornobai/beton.git/internal/http/handlers/user"

	"github.com/anton-chornobai/beton.git/internal/http/middleware"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/service"
	"github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	productService "github.com/anton-chornobai/beton.git/internal/modules/product/application"
	userService "github.com/anton-chornobai/beton.git/internal/modules/user/application"
)

func SetUpRoutes(
	logger *slog.Logger,
	userService *userService.UserService,
	orderService *application.OrderService,
	productService productService.ProductService,
	userContactService service.UserContactService,
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

	orderHandler := handlers.NewOrdersHandler(logger, orderService)

	userContactHandler := handlers.NewUserContactHandler(&userContactService, logger)

	router := http.NewServeMux()
	//AUTH
	router.HandleFunc("POST /signup", authHandler.SignupByEmail)
	router.HandleFunc("POST /login", authHandler.LoginByEmail)
	router.HandleFunc("POST /verify", authHandler.Verify)
	//PROFILE
	router.HandleFunc("GET /profile", userHandler.GetByID)
	//ORDERS
	router.HandleFunc("POST /v1/orders", middleware.GetUsersID(http.HandlerFunc(orderHandler.Create)))
	router.HandleFunc("GET /v1/orders", orderHandler.Get)
	router.Handle("DELETE /v1/orders/{id}", middleware.AdminOnly(userService, http.HandlerFunc(orderHandler.Delete)))
	//PRODUCTS
	router.Handle("GET /v1/products", http.HandlerFunc(productHandler.GetProducts))
	router.Handle("POST /v1/products", middleware.AdminOnly(userHandler.UserService, http.HandlerFunc(productHandler.Add)))
	router.Handle("GET /v1/products/{id}", http.HandlerFunc(productHandler.GetProductByID))
	router.Handle("DELETE /v1/products/{id}", http.HandlerFunc(productHandler.DeleteByID))
	router.Handle("PATCH /v1/products/{id}", http.HandlerFunc(productHandler.Update))
	//CONTACTS
	router.HandleFunc("POST /contacts", userContactHandler.Post)
	router.HandleFunc("DELETE /contacts/{id}", userContactHandler.Delete)


	return middleware.LogMethodInfo(logger, middleware.CorsMiddleware(router))
}
