package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"notification/internal/config"
	"notification/internal/server"
	"notification/internal/transport/http/handler"
	"notification/internal/transport/http/middleware"
	"time"

	"github.com/go-playground/validator/v10"
)

type App struct {
	srv    *http.Server
	logger *slog.Logger
}

func buildHTTP(
	logger *slog.Logger,
	validate *validator.Validate,
) http.Handler {
	h := handler.NewHandler(logger, validate)

	router := handler.NewRouter(h)

	var result http.Handler = router

	result = middleware.ApplicationJsonMiddleware(logger)(result)
	result = middleware.Recovery(logger)(result)
	result = middleware.RequestIDMiddleware(result)

	return result
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	validate := validator.New()

	httpHandler := buildHTTP(
		logger,
		validate,
	)

	serverApp := server.NewServer(httpHandler, cfg)

	return &App{
		srv:    serverApp,
		logger: logger,
	}
}

func (app *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := app.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	app.logger.Info(
		"server started",
		slog.String("addr", app.srv.Addr),
	)

	select {
	case <-ctx.Done():
		app.logger.Info("shutdown signal received")

		return app.shutdown(ctx)

	case err := <-errCh:
		return err
	}
}

func (app *App) shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := app.srv.Shutdown(shutdownCtx); err != nil {
		app.logger.Error(
			"shutdown failed",
			slog.Any("error", err),
		)

		return err
	}

	app.logger.Info("server stopped")

	return nil
}
