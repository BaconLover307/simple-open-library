package middleware

import (
	"net/http"
	"os"
	"simple-open-library/helper"
	"simple-open-library/model/web"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware  {
	return &AuthMiddleware{handler}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	apiKey := os.Getenv("X-API-KEY") 
	if apiKey == request.Header.Get("X-API-KEY") {
		// ok
		 middleware.Handler.ServeHTTP(writer, request)
	} else {
		// error
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code: http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}

		helper.WriteResponseBody(writer, webResponse)
	}
}