package bootstrap

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/anton-chornobai/beton.git/internal/http/routes"

	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	productService "github.com/anton-chornobai/beton.git/internal/modules/product/application"
	productInfra "github.com/anton-chornobai/beton.git/internal/modules/product/infra"

	userContactRepo "github.com/anton-chornobai/beton.git/internal/modules/contact/infra"
	"github.com/anton-chornobai/beton.git/internal/modules/contact/service"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"

	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"
)

func App(db *sql.DB) http.Handler {
	ctx := context.Background()
	//LOGGER
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	//PASSWORD HASHER FOR USERS SERVICE
	passwordHasher := &infra.PasswordHasher{}
	// VERIFICATION 5 DIGIT CODE FOR THE USER SERVICE
	verificationCodeManager := &infra.VerificationaCodeManager{}
	// TOKEN
	tokenManager := jwtmanager.NewTokenService()
	//USER
	userRepo := &infra.UserRepository{DB: db}
	userService := application.NewUserService(userRepo, tokenManager, passwordHasher, log, verificationCodeManager)
	// Google Cloud Platform
	GCPUploader := productInfra.NewGCPUploader(NewGCSClient(ctx))
	//PRODUCTS
	productRepo := &productInfra.ProductRepository{DB: db}
	productService, err := productService.NewProductService(productRepo, GCPUploader, log)
	if err != nil {
		log.Error("failed to create product service", "error", err)
		panic(err)
	}
	//ORDER
	ordersRepo := &ordersRepo.OrdersRepository{DB: db}
	orderService := ordersApp.NewOrderService(ordersRepo, log)
	//CONTACT
	userContactRepo := &userContactRepo.UserContactRepository{DB: db}
	userContactService := service.NewUserContactService(userContactRepo)

	router := routes.SetUpRoutes(log, userService, orderService, *productService, *userContactService)

	return router
}

func NewGCSClient(ctx context.Context) *storage.Client {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create GCS client: %v", err)
	}

	return client
}
