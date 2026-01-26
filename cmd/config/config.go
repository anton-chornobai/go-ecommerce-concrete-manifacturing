package config

import (
	yaml "gopkg.in/yaml.v2"
	"os"
)

type AppConfig struct {
	Env 	string `yaml:"env"`
	Secret string `yaml:"secret"`
	DBPath string `yaml:"db_path"`
	Port   int    `yaml:"port"`
}

type Config struct {
	App AppConfig `yaml:"app"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
