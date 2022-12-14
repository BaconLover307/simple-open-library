package service

import (
	"context"
	"database/sql"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/lib"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"

	"github.com/go-playground/validator/v10"
)

type openLibraryService struct {
	OpenLibraryLib lib.OpenLibraryLib
	DB *sql.DB
	Validate *validator.Validate
}

func NewOpenLibraryService(openLibraryLib lib.OpenLibraryLib, db *sql.DB, validate *validator.Validate) LibraryService {
	return &openLibraryService{
		OpenLibraryLib: openLibraryLib,
		DB: db,
		Validate: validate,
	}
}

func (service openLibraryService) BrowseBySubject(ctx context.Context, request web.SubjectRequest) web.SubjectResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	
	openLibSubjectsResponse := service.OpenLibraryLib.BrowseSubjects(ctx, request.Subject, request.Page)
	if openLibSubjectsResponse.WorkCount == 0 {
		panic(exception.NewNotFoundError("subject not found"))
	}

	var libraryBooks []domain.Book
	for _, book := range openLibSubjectsResponse.Works {
		libraryBooks = append(libraryBooks, domain.NewBookFromOpenLibrary(&book))
	}
	
	return web.SubjectResponse{
		Subject: request.Subject,
		BookCount: openLibSubjectsResponse.WorkCount,
		Page: request.Page,
		Books: web.NewBookResponses(libraryBooks),
	}
}