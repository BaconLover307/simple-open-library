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

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{}
}

func (repo BookRepositoryImpl) SaveBook(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book {
	query := "INSERT INTO book(bookId, title, edition) VALUES(?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, book.BookId, book.Title, book.Edition)
	helper.PanicIfError(err)

	return book

}

// func (repo BookRepositoryImpl) FindBook(ctx context.Context, tx *sql.Tx, book domain.Book) (domain.Book, error) {
// 	query := "SELECT bookId FROM book WHERE title = ? AND author = ? AND edition = ?"
// 	rows, err := tx.QueryContext(ctx, query, book.Title, book.Author, book.Edition)
// 	helper.PanicIfError(err)
// 	defer rows.Close()

// 	resultBook := domain.Book{}
// 	if rows.Next() {
// 		resultBook = book
// 		err = rows.Scan(&resultBook.BookId)
// 		helper.PanicIfError(err)

// 		return resultBook, nil
// 	} else {
// 		return resultBook, exception.NewNotFoundError("book not found") 
// 	}
// }

func (repo BookRepositoryImpl) FindBookById(ctx context.Context, tx *sql.Tx, bookId string) (domain.Book, error) {
	query := `
	SELECT b.bookId, b.title, b.edition, a.authorId, a.name
	FROM author a JOIN authored ab ON a.authorId = ab.authorId
		JOIN book b ON ab.bookId = b.bookId
	WHERE b.bookId = ?;
	`
	rows, err := tx.QueryContext(ctx, query, bookId)
	helper.PanicIfError(err)
	defer rows.Close()

	book := domain.Book{}
	if rows.Next() {
		var authors []domain.Author
		author := domain.Author{}
		err = rows.Scan(&book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
		helper.PanicIfError(err)
		authors = append(authors, author)
		
		// Get remaining authors
		for rows.Next() {
			author := domain.Author{}			
			err = rows.Scan(&book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
			helper.PanicIfError(err)
			authors = append(authors, author)
		}
		book.Authors = authors

		return book, nil
	} else {
		return book, exception.NewNotFoundError("book not found") 
	}
}

func (repo BookRepositoryImpl) Authored(ctx context.Context, tx *sql.Tx, authorId string, bookId string) {
	query := "INSERT INTO authored(authorId, bookId) VALUES(?, ?)"

	_, err := tx.ExecContext(ctx, query, authorId, bookId)
	helper.PanicIfError(err)
}

func (repo BookRepositoryImpl) SaveAuthor(ctx context.Context, tx *sql.Tx, author domain.Author) domain.Author {
	query := "INSERT INTO author(authorId, name) VALUES(?, ?)"

	_, err := tx.ExecContext(ctx, query, author.AuthorId, author.Name)
	helper.PanicIfError(err)

	return author
}

func (repo BookRepositoryImpl) FindAuthor(ctx context.Context, tx *sql.Tx, authorId string) (domain.Author, error) {
	query := "SELECT authorId, name FROM author WHERE authorId = ?"
	rows, err := tx.QueryContext(ctx, query, authorId)
	helper.PanicIfError(err)
	defer rows.Close()

	author := domain.Author{}
	if rows.Next() {
		err = rows.Scan(&author.AuthorId, &author.Name)
		helper.PanicIfError(err)

		return author, nil
	} else {
		return author, exception.NewNotFoundError("author not found") 
	}
}

