package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type LibraryController interface {
	BrowseBySubject(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}