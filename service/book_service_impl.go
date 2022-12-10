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
	DB *sql.DB
	Validate *validator.Validate
}

func NewBookService(bookRepo repository.BookRepository, db *sql.DB, validate *validator.Validate) BookService {
	return &BookServiceImpl{
		BookRepo: bookRepo,
		DB: db,
		Validate: validate,
	}
}

func (service BookServiceImpl) SaveBook(ctx context.Context, request web.BookRequest) web.BookResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book, err := service.BookRepo.FindBookById(ctx, tx, request.BookId)
	if err == nil {
		return web.NewBookResponse(&book)
	}

	book = domain.Book{
		Title: request.Title,
		Edition: request.Edition,
	}
	book = service.BookRepo.SaveBook(ctx, tx, book)
	
	var authors []domain.Author
	for _, authorRequest := range request.Authors {
		author, err := service.BookRepo.FindAuthor(ctx, tx, authorRequest.AuthorId)
		if (err != nil) {
			author = service.BookRepo.SaveAuthor(ctx, tx, author)
			service.BookRepo.Authored(ctx, tx, author.AuthorId, book.BookId)
		}
		authors = append(authors, author)
	}
	book.Authors = authors

	return web.NewBookResponse(&book)
}

// func (service BookServiceImpl) FindBook(ctx context.Context, request web.BookRequest) web.BookResponse {
// 	err := service.Validate.Struct(request)
// 	helper.PanicIfError(err)

// 	tx, err := service.DB.Begin()
// 	helper.PanicIfError(err)
// 	defer helper.CommitOrRollback(tx)

// 	book := domain.Book{
// 		Title: request.Title,
// 		Author: request.Author,
// 		Edition: request.Edition,
// 	}
// 	book, err = service.BookRepo.FindBook(ctx, tx, book)
// 	helper.PanicIfError(err)

// 	return web.NewBookResponse(&book)
// }

func (service BookServiceImpl) FindBookById(ctx context.Context, bookId string) web.BookResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	book, err := service.BookRepo.FindBookById(ctx, tx, bookId)
	helper.PanicIfError(err)

	return web.NewBookResponse(&book)
}

