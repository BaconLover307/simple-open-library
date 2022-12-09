package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type BookRepository interface {
	Save(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book
	FindBook(ctx context.Context, tx *sql.Tx, book domain.Book) (domain.Book, error)
	FindById(ctx context.Context, tx *sql.Tx, bookId int) (domain.Book, error)
}