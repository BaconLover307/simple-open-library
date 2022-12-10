package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type BookController interface {
	ListBooks(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}