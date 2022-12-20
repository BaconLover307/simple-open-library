package app

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func NewServer(mux *chi.Mux) *http.Server {
	return &http.Server{
		Addr:    os.Getenv("SV_HOST") + ":" + os.Getenv("SV_PORT"),
		Handler: mux,
	}
}
