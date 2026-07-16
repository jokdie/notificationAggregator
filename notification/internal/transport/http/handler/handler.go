package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"notification/internal/model"
	"notification/internal/transport/http/response"

	"github.com/go-playground/validator/v10"
)

var errMultipleJSON = errors.New("request body must contain only one JSON object")

type Service interface {
	Send(ctx context.Context, req model.NotificationRequest) (model.NotificationResult, error)
}

type Handler struct {
	logger   *slog.Logger
	validate *validator.Validate
	service  Service
}

func NewHandler(
	logger *slog.Logger,
	validate *validator.Validate,
	service Service,
) *Handler {
	return &Handler{
		logger:   logger,
		validate: validate,
		service:  service,
	}
}

func (h *Handler) decodeRequest(w http.ResponseWriter, r *http.Request) (model.NotificationRequest, error) {
	var req model.NotificationRequest

	const maxBodySize = 1 << 20 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return model.NotificationRequest{}, err
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		return model.NotificationRequest{}, errMultipleJSON
	}

	return req, nil
}

func (h *Handler) handleDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError

	switch {
	case errors.As(err, &maxBytesErr):
		errCode := http.StatusRequestEntityTooLarge

		response.WriteJSON(h.logger, w, errCode, model.ErrorResponse{
			Code:    errCode,
			Message: "Request body too large",
		})

	case errors.Is(err, errMultipleJSON):
		response.BadRequest(h.logger, w, "Only one JSON object is allowed")

	default:
		response.BadRequest(h.logger, w, "Uncorrected JSON")
	}
}
