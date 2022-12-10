package app

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func NewServer(router *httprouter.Router) *http.Server {
	return &http.Server{
		Addr:    os.Getenv("SV_HOST") + ":" + os.Getenv("SV_PORT"),
		Handler: router,
	}
}