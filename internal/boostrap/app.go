package bootstrap

import (
	"database/sql"
	"net/http"

	"github.com/anton-chornobai/beton.git/internal/http/routes"

	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

// import (
// 	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
// 	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
// )

func App(db *sql.DB) *http.ServeMux {
	passwordHasher := &infra.PasswordHasher{}
	tokenManager := jwtmanager.NewTokenService()
	userRepo := &infra.UserRepository{DB: db}
	userService := application.NewUserService(userRepo, tokenManager, passwordHasher)

	ordersRepo := &ordersRepo.OrdersRepository{DB: db}
	orderService := ordersApp.NewOrderService(ordersRepo)

	router := routes.SetUpRoutes(userService, orderService)

	return router
}