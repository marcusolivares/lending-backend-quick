package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	HTTPAddr    string
}

func Load() (Config, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = os.Getenv("PORT")
	}
	if addr == "" {
		addr = ":8080"
	}

	return Config{DatabaseURL: dsn, HTTPAddr: addr}, nil
}
