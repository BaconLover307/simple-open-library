package service

import (
	"context"
	"database/sql"
	"reflect"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"

	"github.com/go-playground/validator/v10"
)

type BookService interface {
	SaveBook(ctx context.Context, request web.BookRequest) web.BookResponse
	FindBookById(ctx context.Context, bookId string) web.BookResponse
	FindAllBooks(ctx context.Context) []web.BookResponse
}

type bookService struct {
	BookRepo repository.BookRepository
	DB *sql.DB
	Validate *validator.Validate
}

func NewBookService(bookRepo repository.BookRepository, db *sql.DB, validate *validator.Validate) BookService {
	return &bookService{
		BookRepo: bookRepo,
		DB: db,
		Validate: validate,
	}
}

func (service bookService) SaveBook(ctx context.Context, request web.BookRequest) web.BookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book, err := service.BookRepo.FindBookById(ctx, tx, request.BookId)
	if (err == nil && !reflect.DeepEqual(web.NewBook(&request), book)) {
		panic(exception.NewConflictError("cannot overwrite existing book. please insert correct book data"))
	}
	if err == nil {
		return web.NewBookResponse(&book)
	}

	book = domain.Book{
		BookId: request.BookId,
		Title: request.Title,
		Edition: request.Edition,
	}
	book = service.BookRepo.SaveBook(ctx, tx, book)
	
	var authors []domain.Author
	for _, authorRequest := range request.Authors {
		author, err := service.BookRepo.FindAuthor(ctx, tx, authorRequest.AuthorId)
		if (err != nil) {
			newAuthor := domain.Author{
				AuthorId: authorRequest.AuthorId,
				Name: authorRequest.Name,
			}
			author = service.BookRepo.SaveAuthor(ctx, tx, newAuthor)
		}
		service.BookRepo.Authored(ctx, tx, author.AuthorId, book.BookId)
		authors = append(authors, author)
	}
	book.Authors = authors

	return web.NewBookResponse(&book)
}

func (service bookService) FindBookById(ctx context.Context, bookId string) web.BookResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book, err := service.BookRepo.FindBookById(ctx, tx, bookId)
	helper.PanicIfError(err)

	return web.NewBookResponse(&book)
}

func (service bookService) FindAllBooks(ctx context.Context) []web.BookResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickups := service.BookRepo.FindAllBooks(ctx, tx)

	return web.NewBookResponses(pickups)
}
