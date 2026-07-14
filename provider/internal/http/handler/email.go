package handler

import (
	"net/http"
	"provider/internal/http/response"
	"provider/internal/model"
)

func (h *Handler) email(w http.ResponseWriter, r *http.Request) {
	req, err := h.decodeRequest(w, r)

	if err != nil {
		response.BadRequest(h.logger, w)

		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(h.logger, w)

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
