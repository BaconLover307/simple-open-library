package service

import (
	"context"
	"simple-open-library/model/web"
)

type BookService interface {
	BrowseSubject(ctx context.Context, request web.SubjectRequest) []web.BookResponse
	Save(ctx context.Context, request web.BookRequest) web.BookResponse
	FindBook(ctx context.Context, request web.BookRequest) web.BookResponse
	FindById(ctx context.Context, bookId int) web.BookResponse
}