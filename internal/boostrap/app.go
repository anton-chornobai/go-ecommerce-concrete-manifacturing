package bootstrap

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/anton-chornobai/beton.git/internal/http/routes"

	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	productService "github.com/anton-chornobai/beton.git/internal/modules/product/application"
	productRepo "github.com/anton-chornobai/beton.git/internal/modules/product/infra"

	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

// import (
// 	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
// 	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
// )

func App(db *sql.DB) http.Handler{
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	passwordHasher := &infra.PasswordHasher{}
	tokenManager := jwtmanager.NewTokenService()
	userRepo := &infra.UserRepository{DB: db}
	userService := application.NewUserService(userRepo, tokenManager, passwordHasher, log)

	productRepo := &productRepo.ProductRepository{DB: db}
	productService, err := productService.NewProductService(productRepo)
	if err != nil {
		
	}

	ordersRepo := &ordersRepo.OrdersRepository{DB: db}
	orderService := ordersApp.NewOrderService(ordersRepo)

	router := routes.SetUpRoutes(userService, orderService, *productService)

	return router
}