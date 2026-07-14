package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"provider/internal/http/handler"
	"provider/internal/http/middleware"
	"provider/internal/provider/email"
	"provider/internal/server"
	"provider/internal/service"
	"time"

	"github.com/go-playground/validator/v10"
)

type App struct {
	srv    *http.Server
	logger *slog.Logger
}

func New(logger *slog.Logger) *App {
	validate := validator.New()

	emailProvider := email.New(logger)
	notificationService := service.NewNotificationService(emailProvider)

	handlerApp := handler.NewHandler(validate, logger, notificationService)
	routerApp := handler.NewRouter(handlerApp)

	var h http.Handler = routerApp

	h = middleware.ApplicationJsonMiddleware(logger)(h)
	h = middleware.RequestIDMiddleware(h)

	serverApp := server.NewServer(h)

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

		return app.shutdown(context.Background())

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
