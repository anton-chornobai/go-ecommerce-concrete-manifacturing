package main

import (
	"fmt"
	"github.com/anton-chornobai/beton.git/cmd/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
	"github.com/anton-chornobai/beton.git/internal/routes"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load("../../.env")
	
	if err != nil {
		log.Fatal(err)

	}

	cfg, err := config.LoadConfig("../../configs/app.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(os.Environ())
	fmt.Println("SECRET =", os.Getenv("SECRET"))

	logger := setupLogger(cfg.App.Env)

	log.Println("SQLite DB path:", cfg.App.DBPath)
	db := db.Connect(cfg.App.DBPath)

	defer db.Close()

	userRepo := &infra.UserRepository{DB: db, Logger: logger}
	userAppService := application.NewUserService(userRepo, logger)

	ordersRepo := &ordersRepo.OrdersRepository{DB: db}
	orderService := ordersApp.NewOrderService(ordersRepo)

	handler := routes.SetUpRouter(userAppService, orderService)

	myService := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: handler,
	}
	fmt.Printf("Server is running on port: http://localhost%s\n", myService.Addr)
	if err := myService.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
