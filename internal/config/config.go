package config

import (
	"errors"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Env    string   `yaml:"env"`
	DB     DBConfig `yaml:"database"`
	Port   int      `yaml:"port"`
}

type DBConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	SSLmode  string `yaml:"sslmode"`
	Password int    `yaml:"password"`
	Port     int    `yaml:"port"`
}

func LoadConfig() (*Config, error) {
	cfg_path := os.Getenv("CONFIG_PATH")
	if cfg_path == "" {
		return nil, errors.New("no cfg file found")
	}

	data, err := os.ReadFile(cfg_path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
