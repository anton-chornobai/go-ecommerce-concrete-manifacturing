package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath string
	Port   string
}

func LoadConfig() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		DBPath: os.Getenv("DB_PATH"),
		Port:   ":" + os.Getenv("PORT"),
	}
}
