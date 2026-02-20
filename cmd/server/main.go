package main

import (
	"fmt"

	"log"
	"net/http"
	"strconv"

	"github.com/anton-chornobai/beton.git/internal/boostrap"
	"github.com/anton-chornobai/beton.git/internal/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	"github.com/joho/godotenv"
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

	// logger := config.SetupLogger(cfg.App.Env)

	log.Println("SQLite DB path:", cfg.App.DBPath)
	db := db.Connect(cfg.App.DBPath)

	defer db.Close()

	router := bootstrap.App(db)

	myService := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: router,
	}
	fmt.Printf("Server is running on port: http://localhost%s\n", myService.Addr)
	if err := myService.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}
