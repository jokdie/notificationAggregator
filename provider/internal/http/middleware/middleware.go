package middleware

import (
	"context"
	"log/slog"
	"mime"
	"net/http"
	"provider/internal/http/response"
	"provider/internal/model"
	"provider/internal/requestid"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-GUID")

		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set("X-GUID", requestID)

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
