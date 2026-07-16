package middleware

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"notification/internal/model"
	"notification/internal/requestid"
	"notification/internal/transport/http/response"
	"runtime/debug"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestid.Key, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ApplicationJsonMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			if err != nil || mediaType != "application/json" {
				errCode := http.StatusUnsupportedMediaType
				response.WriteJSON(logger, w, errCode, model.ErrorResponse{Code: errCode, Message: "Unsupported Media Type"})

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(
						"panic recovered",
						slog.Any("panic", rec),
						slog.String("stack", string(debug.Stack())),
						slog.String("request_id", requestid.Get(r.Context())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.RequestURI()),
					)

					response.Internal(logger, w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
