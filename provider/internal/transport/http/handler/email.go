package handler

import (
	"log/slog"
	"net/http"
	"provider/internal/model"
	"provider/internal/requestid"
	"provider/internal/transport/http/response"
)

func (h *Handler) email(w http.ResponseWriter, r *http.Request) {
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
	)

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(h.logger, w, "Validation Error")

		return
	}

	if err := h.notificationService.Send(
		r.Context(),
		model.Email,
		req,
	); err != nil {
		response.Internal(h.logger, w)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
