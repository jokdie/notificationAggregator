package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h *Handler) notifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := struct{}{}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("json encode error:", slog.Any("error", err))
	}
}
