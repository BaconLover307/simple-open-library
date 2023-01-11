package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LibraryController interface {
	BrowseBySubject(writer http.ResponseWriter, request *http.Request)
}

type libraryController struct {
	LibraryService service.LibraryService
}

func NewLibraryController(libraryService service.LibraryService) LibraryController {
	return &libraryController{
		LibraryService: libraryService,
	}
}

func (controller libraryController) BrowseBySubject(writer http.ResponseWriter, request *http.Request) {
	subject := chi.URLParam(request, "subject")
	pageQuery := request.URL.Query().Get("page")
	if pageQuery == "" {
		pageQuery = "1"
	}
	page, err := strconv.Atoi(pageQuery)
	helper.PanicIfError(err)

	subjectResponse := controller.LibraryService.BrowseBySubject(request.Context(), web.SubjectRequest{Subject: subject, Page: page})
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   subjectResponse,
	}

	helper.WriteResponseBody(writer, webResponse)
}
