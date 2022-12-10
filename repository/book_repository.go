package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type BookRepository interface {
	SaveBook(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book
	FindBookById(ctx context.Context, tx *sql.Tx, bookId string) (domain.Book, error)
	FindAllBooks(ctx context.Context, tx *sql.Tx) []domain.Book 
	Authored(ctx context.Context, tx *sql.Tx, authorId string, bookId string)
	SaveAuthor(ctx context.Context, tx *sql.Tx, author domain.Author) domain.Author
	FindAuthor(ctx context.Context, tx *sql.Tx, authorId string) (domain.Author, error)
}