package service

import (
	"context"
	"database/sql"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/repository"

	"github.com/go-playground/validator/v10"
)

type LibraryServiceImpl struct {
	LibraryRepo repository.LibraryRepository
	DB *sql.DB
	Validate *validator.Validate
}

func NewLibraryService(libraryRepo repository.LibraryRepository, db *sql.DB, validate *validator.Validate) LibraryService {
	return &LibraryServiceImpl{
		LibraryRepo: libraryRepo,
		DB: db,
		Validate: validate,
	}
}

func (service LibraryServiceImpl) BrowseBySubject(ctx context.Context, request web.SubjectRequest) []web.BookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	
	books, err := service.LibraryRepo.Subjects(ctx, request.Subject, request.Page)
	helper.PanicIfError(err)
	
	return web.NewBookResponses(books)
}