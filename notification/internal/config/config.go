package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func Load() *Config {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	return &Config{
		HTTPAddr:          ":" + port,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
