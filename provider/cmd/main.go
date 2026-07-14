package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"provider/internal/app"
	"provider/internal/config"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	application := app.New(cfg, logger)

	if err := application.Run(ctx); err != nil {
		logger.Error("application exited", slog.Any("error", err))
	}
}
