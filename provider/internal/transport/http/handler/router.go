package handler

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /provider/v1/email", h.email)
	mux.HandleFunc("POST /provider/v1/sms", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /provider/v1/push", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
