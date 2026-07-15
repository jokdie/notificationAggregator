package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"provider/internal/config"
	"provider/internal/model"
	"provider/internal/provider"
	"provider/internal/provider/email"
	"provider/internal/server"
	"provider/internal/service"
	"provider/internal/transport/http/handler"
	"provider/internal/transport/http/middleware"
	"time"

	"github.com/go-playground/validator/v10"
)

type App struct {
	srv    *http.Server
	logger *slog.Logger
}

func buildProviders(logger *slog.Logger) map[model.Channel]provider.Provider {
	return map[model.Channel]provider.Provider{
		model.Email: email.New(logger),
	}
}

func buildHTTP(
	logger *slog.Logger,
	validate *validator.Validate,
	notificationService *service.NotificationService,
) http.Handler {
	h := handler.NewHandler(
		validate,
		logger,
		notificationService,
	)

	router := handler.NewRouter(h)

	var result http.Handler = router

	result = middleware.ApplicationJsonMiddleware(logger)(result)
	result = middleware.Recovery(logger)(result)
	result = middleware.RequestIDMiddleware(result)

	return result
}

func New(cfg *config.Config, logger *slog.Logger) *App {
	validate := validator.New()

	providers := buildProviders(logger)
	notificationService := service.NewNotificationService(providers)

	httpHandler := buildHTTP(
		logger,
		validate,
		notificationService,
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

	// Здесь также закрываем соединения с БД и фоновые задачи (например, db.Close())

	app.logger.Info("server stopped")

	return nil
}
