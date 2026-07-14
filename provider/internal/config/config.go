package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr          string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return &Config{
		HTTPAddr:          ":" + port,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
