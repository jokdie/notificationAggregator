package main

import (
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {})
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return
	}
}
