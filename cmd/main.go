package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/anton-chornobai/beton.git/internal/boostrap"
	"github.com/anton-chornobai/beton.git/internal/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	"github.com/joho/godotenv"
)

//go:embed swagger.html
var swaggerUI embed.FS

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("no .env file found, using environment variables")
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	connStr := db.GetDBConnStr(config.DB)
	db, err := db.OpenPostgre(connStr)
	if err != nil {
		log.Fatalf("failed to open db %v", err)
	}
	defer db.Close()
	fmt.Printf("Opened DB connection on port: %d \n", config.DB.Port)

	router := bootstrap.App(db)

	myService := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}
	fmt.Printf("Server is running on port: %s\n", myService.Addr)
	if err := myService.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}
