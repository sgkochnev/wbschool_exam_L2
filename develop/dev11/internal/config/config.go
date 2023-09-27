package config

import (
	"os"
)

// Config - config
type Config struct {
	Address     string
	StroagePath string
}

// New create new config from env variables
func New() *Config {
	cfg := &Config{}
	cfg.Address = os.Getenv("HTTP_SERVER_ADDRESS")
	cfg.StroagePath = os.Getenv("STORAGE_PATH")

	if cfg.Address == "" {
		cfg.Address = "localhost:8080"
	}

	if cfg.StroagePath == "" {
		cfg.StroagePath = "./storage"
	}

	return cfg
}
