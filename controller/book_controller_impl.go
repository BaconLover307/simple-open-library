package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"

	"github.com/julienschmidt/httprouter"
)

type BookControllerImpl struct {
	BookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &BookControllerImpl{BookService: bookService}
}

func (controller BookControllerImpl) ListBooks(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bookResponses := controller.BookService.FindAllBooks(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   bookResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}
