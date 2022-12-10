package service

import (
	"context"
	"simple-open-library/model/web"
)

type BookService interface {
	BrowseSubject(ctx context.Context, request web.SubjectRequest) []web.BookResponse
	SaveBook(ctx context.Context, request web.BookRequest) web.BookResponse
	FindBookById(ctx context.Context, bookId string) web.BookResponse
}