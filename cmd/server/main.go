package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"github.com/anton-chornobai/beton.git/cmd/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	"github.com/anton-chornobai/beton.git/internal/modules/user/application"
	"github.com/anton-chornobai/beton.git/internal/modules/user/domain"
	"github.com/anton-chornobai/beton.git/internal/modules/user/infra"
	"github.com/anton-chornobai/beton.git/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	dbPath, err := filepath.Abs(cfg.DBPath)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("SQLite DB path:", dbPath)
	db := db.Connect(dbPath)
	
	defer db.Close()

	userRepo := infra.UserRepository{DB: db}
	userDomainServices := domain.NewService(&userRepo)
	userAppService := application.NewUserService(userDomainServices)

	handler := routes.SetUpRouter(userAppService)

	myService := &http.Server{
		Addr:   cfg.Port,
		Handler: handler,
	}
	fmt.Printf("Server is running on port: http://localhost%s\n", myService.Addr)
	if err := myService.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}

}
