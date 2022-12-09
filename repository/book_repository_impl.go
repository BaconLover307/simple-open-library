package repository

import (
	"context"
	"database/sql"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
)

type BookRepositoryImpl struct {
}

func (repo BookRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book {
	query := "INSERT INTO book(title, author, edition) VALUES(?, ?, ?)"

	result, err := tx.ExecContext(ctx, query, book.Title, book.Author, book.Edition)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	book.BookId = int(id)
	return book

}

func (repo BookRepositoryImpl) FindBook(ctx context.Context, tx *sql.Tx, book domain.Book) (domain.Book, error) {
	query := "SELECT bookId FROM book WHERE title = ? AND author = ? AND edition = ?"
	rows, err := tx.QueryContext(ctx, query, book.Title, book.Author, book.Edition)
	helper.PanicIfError(err)
	defer rows.Close()

	resultBook := domain.Book{}
	if rows.Next() {
		resultBook = book
		err = rows.Scan(&resultBook.BookId)
		helper.PanicIfError(err)

		return resultBook, nil
	} else {
		return resultBook, exception.NewNotFoundError("book not found") 
	}
}

func (repo BookRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, bookId int) (domain.Book, error) {
	query := "SELECT bookId, title, author, edition FROM book WHERE bookId = ?"
	rows, err := tx.QueryContext(ctx, query, bookId)
	helper.PanicIfError(err)
	defer rows.Close()

	book := domain.Book{}
	if rows.Next() {
		err = rows.Scan(&book.BookId, &book.Title, &book.Author, &book.Edition)
		helper.PanicIfError(err)

		return book, nil
	} else {
		return book, exception.NewNotFoundError("book not found") 
	}
}

