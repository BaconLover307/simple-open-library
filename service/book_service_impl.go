package service

import (
	"context"
	"database/sql"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"

	"github.com/go-playground/validator/v10"
)

type BookServiceImpl struct {
	BookRepo repository.BookRepository
	LibraryRepo repository.LibraryRepository
	DB *sql.DB
	Validate *validator.Validate
}

func NewBookService(bookRepo repository.BookRepository, libraryRepo repository.LibraryRepository, db *sql.DB, validate *validator.Validate) BookService {
	return &BookServiceImpl{
		BookRepo: bookRepo,
		LibraryRepo: libraryRepo,
		DB: db,
		Validate: validate,
	}
}

func (service BookServiceImpl) BrowseSubject(ctx context.Context, request web.SubjectRequest) []web.LibraryBookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)
	
	libraryBooks, err := service.LibraryRepo.Subjects(ctx, request.Subject, request.Page)
	helper.PanicIfError(err)
	
	return web.NewLibraryBookResponses(libraryBooks)
	
}

func (service BookServiceImpl) Save(ctx context.Context, request web.BookRequest) web.BookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book := domain.Book{
		Title: request.Title,
		Author: request.Author,
		Edition: request.Edition,
	}
	book = service.BookRepo.Save(ctx, tx, book)

	return web.NewBookResponse(&book)
}

func (service BookServiceImpl) FindBook(ctx context.Context, request web.BookRequest) web.BookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book := domain.Book{
		Title: request.Title,
		Author: request.Author,
		Edition: request.Edition,
	}
	book, err = service.BookRepo.FindBook(ctx, tx, book)
	helper.PanicIfError(err)

	return web.NewBookResponse(&book)
}

func (service BookServiceImpl) FindById(ctx context.Context, bookId int) web.BookResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book, err := service.BookRepo.FindById(ctx, tx, bookId)
	helper.PanicIfError(err)

	return web.NewBookResponse(&book)
}

