package handler

import (
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /provider/email", h.email)
	mux.HandleFunc("POST /provider/sms", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /provider/push", func(w http.ResponseWriter, r *http.Request) {})

	return mux
}
