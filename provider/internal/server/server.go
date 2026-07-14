package server

import "net/http"

func NewServer(routerApp http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routerApp,
	}

	return srv
}
