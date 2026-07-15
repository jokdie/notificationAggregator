package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"provider/internal/model"
)

func WriteJSON(logger *slog.Logger, w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("json encode error:", slog.Any("error", err))
	}
}

func BadRequest(logger *slog.Logger, w http.ResponseWriter, message string) {
	errCode := http.StatusBadRequest
	WriteJSON(logger, w, errCode, model.ErrorResponse{Code: errCode, Message: message})
}

func Internal(logger *slog.Logger, w http.ResponseWriter) {
	errCode := http.StatusInternalServerError
	WriteJSON(logger, w, errCode, model.ErrorResponse{Code: errCode, Message: "Internal server error"})
}
