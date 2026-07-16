package handler

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /notifications/v1/send", h.notifications)

	return mux
}
