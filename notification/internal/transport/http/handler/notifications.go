package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"notification/internal/requestid"
	"notification/internal/transport/http/response"
)

func (h *Handler) notifications(w http.ResponseWriter, r *http.Request) {
	req, err := h.decodeRequest(w, r)

	if err != nil {
		h.handleDecodeError(w, err)

		return
	}

	h.logger.Info(
		"notification request",
		slog.String("request_id", requestid.Get(r.Context())),
		slog.String("method", r.Method),
		slog.String("path", r.URL.RequestURI()),
		slog.Int("user_id", req.UserID),
		slog.Any("channels", req.Channels),
	)

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(h.logger, w, "Validation Error")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, err := h.service.Send(r.Context(), req)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		h.logger.Error("json encode error:", slog.Any("error", err))
	}
}
