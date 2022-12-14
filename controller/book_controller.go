package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"

	"github.com/julienschmidt/httprouter"
)

type BookController interface {
	ListBooks(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type bookController struct {
	BookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &bookController{BookService: bookService}
}

func (controller bookController) ListBooks(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bookResponses := controller.BookService.FindAllBooks(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   bookResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}