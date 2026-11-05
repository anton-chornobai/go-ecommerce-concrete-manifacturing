package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/anton-chornobai/beton.git/internal/config"
	"github.com/anton-chornobai/beton.git/internal/http/routes"

	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	productService "github.com/anton-chornobai/beton.git/internal/modules/product/application"
	productRepo "github.com/anton-chornobai/beton.git/internal/modules/product/infra/postgres"

	"github.com/anton-chornobai/beton.git/internal/modules/contact/application"
	userContactRepo "github.com/anton-chornobai/beton.git/internal/modules/contact/infra"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"

	jwtmanager "github.com/anton-chornobai/beton.git/internal/modules/user/infra/jwt"

	s3storage "github.com/anton-chornobai/beton.git/internal/modules/product/infra/image_storage"
	awsservice "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func App(db *sql.DB, cfg *config.Config) http.Handler {
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

	// AWS S3 Image Storage
	imageStorage, err := NewS3Storage(ctx, cfg)
	if err != nil {
		log.Error("failed to initialize S3 storage", "error", err)
		panic(err)
	}
	productRepo := &productRepo.ProductRepository{DB: db}
	productService, err := productService.NewProductService(productRepo, imageStorage, log)
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

	router := routes.SetUpRoutes(log, userService, orderService, productService, *userContactService)

	return router
}

func NewS3Storage(ctx context.Context, cfg *config.Config) (*s3storage.S3Storage, error) {
	awsCfg, err := awsservice.LoadDefaultConfig(
		ctx, awsservice.WithRegion(cfg.S3.Region),
	)

	if err != nil {
		return nil, fmt.Errorf("load AWS configuration: %w", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return s3storage.NewS3Storage(client, cfg.S3.Bucket, cfg.S3.CloudFrontURL), nil
}
