package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"provider/internal/http/response"
	"provider/internal/model"
	"provider/internal/service"

	"github.com/go-playground/validator/v10"
)

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

// @todo удалить все writeJSON из метода decodeRequest тк слишком много ответственностей внутри одного метода
func (h *Handler) decodeRequest(w http.ResponseWriter, r *http.Request) (model.ProviderRequest, error) {
	var req model.ProviderRequest

	const maxBodySize = 1 << 20 // 1 MB

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError

		switch {
		case errors.As(err, &maxBytesErr):
			errCode := http.StatusRequestEntityTooLarge

			response.WriteJSON(h.logger, w, errCode, model.ErrorResponse{
				Code:    errCode,
				Message: "Request body too large",
			})

		default:
			errCode := http.StatusBadRequest

			response.WriteJSON(h.logger, w, errCode, model.ErrorResponse{
				Code:    errCode,
				Message: "Некорректный JSON",
			})
		}

		return model.ProviderRequest{}, err
	}

	if decoder.Decode(&struct{}{}) != io.EOF {
		errCode := http.StatusBadRequest

		response.WriteJSON(h.logger, w, errCode, model.ErrorResponse{
			Code:    errCode,
			Message: "Разрешен только один объект JSON",
		})

		return model.ProviderRequest{}, errors.New("request body must contain only one JSON object")
	}

	return req, nil
}
