package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/anton-chornobai/beton.git/cmd/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	ordersApp "github.com/anton-chornobai/beton.git/internal/modules/orders/application"
	ordersRepo "github.com/anton-chornobai/beton.git/internal/modules/orders/infra"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
	"github.com/anton-chornobai/beton.git/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig("../../configs/app.yaml")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("SQLite DB path:", cfg.App.DBPath)
	db := db.Connect(cfg.App.DBPath)

	defer db.Close()

	userRepo := infra.UserRepository{DB: db}
	userDomainServices := domain.NewService(&userRepo)
	userAppService := application.NewUserService(userDomainServices)

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
