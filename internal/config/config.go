package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Env  string
	Port int
	DB   DBConfig
	S3   S3Config
}

type DBConfig struct {
	Name     string
	Host     string
	User     string
	Password string
	SSLMode  string
	Port     int
}

type S3Config struct {
	Region        string
	Bucket        string
	CloudFrontURL string
}

func LoadConfig() (*Config, error) {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}

	cfg := &Config{
		Env:  os.Getenv("APP_ENV"),
		Port: httpPort,
		DB: DBConfig{
			Name:     os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOSTNAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			Port:     dbPort,
		},
		S3: S3Config{
			Bucket:        os.Getenv("S3_BUCKET_NAME"),
			Region:        os.Getenv("AWS_REGION"),
			CloudFrontURL: os.Getenv("CLOUDFRONT_URL"),
		},
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("HTTP_PORT must be between 1 and 65535")
	}

	if cfg.DB.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if cfg.DB.Host == "" {
		return nil, fmt.Errorf("DB_HOSTNAME is required")
	}
	if cfg.DB.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}
	if cfg.DB.Port < 1 || cfg.DB.Port > 65535 {
		return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
	}

	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}

	if cfg.S3.Bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required")
	}

	if cfg.S3.Region == "" {
		return nil, fmt.Errorf("AWS_REGION is required")
	}

	if cfg.S3.CloudFrontURL == "" {
		return nil, fmt.Errorf("CLOUDFRONT_URL is required")
	}

	return cfg, nil
}
