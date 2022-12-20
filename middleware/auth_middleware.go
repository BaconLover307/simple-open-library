package middleware

import (
	"net/http"
	"os"
	"simple-open-library/helper"
	"simple-open-library/model/route"
	"simple-open-library/model/web"
)

type AuthMiddleware struct {
	Handler  http.Handler
	Excludes *route.Prefixes
}

func NewAuthMiddleware(handler http.Handler, excludes *route.Prefixes) *AuthMiddleware {
	return &AuthMiddleware{
		Handler:  handler,
		Excludes: excludes,
	}
}

func (middleware *AuthMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	apiKey := os.Getenv("X-API-KEY")
	if middleware.Excludes.ContainsPrefix(request.RequestURI) || apiKey == request.Header.Get("X-API-KEY") {
		// ok
		middleware.Handler.ServeHTTP(writer, request)
	} else {
		// error
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}

		helper.WriteResponseBody(writer, webResponse)
	}
}

func ChiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("X-API-KEY")
		if apiKey == r.Header.Get("X-API-KEY") {
			// ok
			next.ServeHTTP(w, r)
		} else {
			// error
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)

			webResponse := web.WebResponse{
				Code:   http.StatusUnauthorized,
				Status: "UNAUTHORIZED",
			}

			helper.WriteResponseBody(w, webResponse)
		}
	})
}
