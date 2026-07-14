package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"provider/internal/app"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	application := app.New(logger)

	if err := application.Run(ctx); err != nil {
		logger.Error("application exited", slog.Any("error", err))
	}
}
