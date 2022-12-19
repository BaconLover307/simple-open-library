package repository

import (
	"context"
	"database/sql"
	"simple-open-library/exception"
	"simple-open-library/helper"
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

type bookRepo struct {
}

func NewBookRepository() BookRepository {
	return &bookRepo{}
}

func (repo bookRepo) SaveBook(ctx context.Context, tx *sql.Tx, book domain.Book) domain.Book {
	query := "INSERT INTO book(book_id, title, edition) VALUES(?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, book.BookId, book.Title, book.Edition)
	helper.PanicIfError(err)

	return book

}

func (repo bookRepo) FindAllBooks(ctx context.Context, tx *sql.Tx) []domain.Book {
	query := `
	SELECT b.book_id, b.title, b.edition, a.author_id, a.name
	FROM author a JOIN authored ab ON a.author_id = ab.author_id
		JOIN book b ON ab.book_id = b.book_id
		ORDER BY b.book_id, a.author_id
	`

	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var books []domain.Book
	book := domain.Book{}
	author := domain.Author{}
	var authors []domain.Author

	for rows.Next() {
		err = rows.Scan(&book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
		helper.PanicIfError(err)
		if len(books) == 0 || books[len(books)-1].BookId != book.BookId {
			authors = nil
			authors = append(authors, author)
			book.Authors = authors
			books = append(books, book)
		} else {
			authors = append(authors, author)
			book.Authors = authors
			books[len(books)-1] = book
		}
	}
	return books
}

func (repo bookRepo) FindBookById(ctx context.Context, tx *sql.Tx, bookId string) (domain.Book, error) {
	query := `
	SELECT b.book_id, b.title, b.edition, a.author_id, a.name
	FROM author a JOIN authored ab ON a.author_id = ab.author_id
		JOIN book b ON ab.book_id = b.book_id
	WHERE b.book_id = ?
	ORDER BY a.author_id
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

func (repo bookRepo) Authored(ctx context.Context, tx *sql.Tx, authorId string, bookId string) {
	query := "INSERT INTO authored(author_id, book_id) VALUES(?, ?)"

	_, err := tx.ExecContext(ctx, query, authorId, bookId)
	helper.PanicIfError(err)
}

func (repo bookRepo) SaveAuthor(ctx context.Context, tx *sql.Tx, author domain.Author) domain.Author {
	query := "INSERT INTO author(author_id, name) VALUES(?, ?)"

	_, err := tx.ExecContext(ctx, query, author.AuthorId, author.Name)
	helper.PanicIfError(err)

	return author
}

func (repo bookRepo) FindAuthor(ctx context.Context, tx *sql.Tx, authorId string) (domain.Author, error) {
	query := "SELECT author_id, name FROM author WHERE author_id = ?"
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
