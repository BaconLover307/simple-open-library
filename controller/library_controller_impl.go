package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type LibraryControllerImpl struct {
	LibraryService service.LibraryService 
}

func NewLibraryController(libraryService service.LibraryService) LibraryController {
	return &LibraryControllerImpl{
		LibraryService: libraryService,
	}
}

func (controller LibraryControllerImpl) BrowseBySubject(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	subject := params.ByName("subject")
	pageQuery := request.URL.Query().Get("page")
	if (pageQuery == "") {
		pageQuery = "1"
	}
	page, err := strconv.Atoi(pageQuery)
	helper.PanicIfError(err)

	controller.LibraryService.BrowseBySubject(request.Context(), web.SubjectRequest{Subject: subject, Page: page})
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}
	
	helper.WriteResponseBody(writer, webResponse)
}