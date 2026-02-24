package main

import (
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/anton-chornobai/beton.git/internal/boostrap"
	"github.com/anton-chornobai/beton.git/internal/config"
	"github.com/anton-chornobai/beton.git/internal/db"
	"github.com/anton-chornobai/beton.git/internal/mail"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// logger := config.SetupLogger(cfg.App.Env)

	connStr := db.GetDBConnStr(config.DB)
	db, err := db.OpenPostgre(connStr)
	if err != nil {
		log.Fatalf("failed to open db %v", err)
	}
	defer db.Close()
	
	mail.SendEmailSample()
	router := bootstrap.App(db)

	myService := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}
	fmt.Printf("Server is running on port: http://localhost%s\n", myService.Addr)
	if err := myService.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed %v", err)
	}
}

