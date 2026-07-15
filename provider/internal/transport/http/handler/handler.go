package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"provider/internal/model"
	"provider/internal/service"
	"provider/internal/transport/http/response"

	"github.com/go-playground/validator/v10"
)

var errMultipleJSON = errors.New("request body must contain only one JSON object")

type Handler struct {
	validate            *validator.Validate
	logger              *slog.Logger
	notificationService *service.NotificationService
}

func NewHandler(
	validate *validator.Validate,
	logger *slog.Logger,
	notificationService *service.NotificationService,
) *Handler {
	return &Handler{
		validate:            validate,
		logger:              logger,
		notificationService: notificationService,
	}
}

func (h *Handler) decodeRequest(w http.ResponseWriter, r *http.Request) (model.ProviderRequest, error) {
	var req model.ProviderRequest

	const maxBodySize = 1 << 20 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return model.ProviderRequest{}, err
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		return model.ProviderRequest{}, errMultipleJSON
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
