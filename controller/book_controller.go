package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
)

type BookController interface {
	ListBooks(writer http.ResponseWriter, request *http.Request)
}

type bookController struct {
	BookService service.BookService
}

func NewBookController(bookService service.BookService) BookController {
	return &bookController{BookService: bookService}
}

func (controller bookController) ListBooks(writer http.ResponseWriter, request *http.Request) {
	bookResponses := controller.BookService.FindAllBooks(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   bookResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}
