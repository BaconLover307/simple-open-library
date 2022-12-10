package app

import (
	"net/http"
	"os"
	"simple-open-library/middleware"
)

func NewServer(authMiddleware *middleware.AuthMiddleware) *http.Server {
	return &http.Server{
		Addr:    os.Getenv("SV_HOST") + ":" + os.Getenv("SV_PORT"),
		Handler: authMiddleware,
	}
}