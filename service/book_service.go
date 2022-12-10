package service

import (
	"context"
	"simple-open-library/model/web"
)

type BookService interface {
	SaveBook(ctx context.Context, request web.BookRequest) web.BookResponse
	FindBookById(ctx context.Context, bookId string) web.BookResponse
	FindAllBooks(ctx context.Context) []web.BookResponse
}