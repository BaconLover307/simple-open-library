package repository_test

import (
	"context"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	inputAuthor1 = domain.Author{
		AuthorId: "ta001",
		Name:     "Jk Rolling",
	}
	inputAuthor2 = domain.Author{
		AuthorId: "ta002",
		Name:     "Mark Mansion",
	}
	inputAuthors1 = []domain.Author{inputAuthor1}
	inputAuthors2 = []domain.Author{inputAuthor1, inputAuthor2}

	inputBook1 = domain.Book{
		BookId:  "tb001",
		Title:   "Test Book",
		Edition: 1,
		Authors: inputAuthors1,
	}
	inputBook2 = domain.Book{
		BookId:  "tb002",
		Title:   "Help Book",
		Edition: 2,
		Authors: inputAuthors2,
	}
	bookColumns       = []string{"bookId", "title", "edition"}
	authorColumns     = []string{"authorId", "name"}
	selectBookColumns = []string{"bookId", "title", "edition", "authorId", "name"}
)

func TestRepoBookSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO book").
		WithArgs(inputBook1.BookId, inputBook1.Title, inputBook1.Edition).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	ctx := context.Background()

	bookRepo := repository.NewBookRepository()
	savedBook := bookRepo.SaveBook(ctx, tx, inputBook1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
	require.Equal(t, inputBook1, savedBook)
}

func TestRepoBookSaveAuthor(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO author").
		WithArgs(inputAuthor1.AuthorId, inputAuthor1.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepoBookAuthored(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO authored").
		WithArgs(inputAuthor1.AuthorId, inputBook1.BookId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepoBookFindAuthorSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author").
		WithArgs(inputAuthor1.AuthorId).
		WillReturnRows(sqlmock.NewRows(authorColumns).AddRow(inputAuthor1.AuthorId, inputAuthor1.Name))

	tx, err := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	authorResult, err := bookRepo.FindAuthor(ctx, tx, inputAuthor1.AuthorId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
	require.Equal(t, inputAuthor1, authorResult)
}

func TestRepoBookFindAuthorFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author").
		WithArgs(inputAuthor1.AuthorId).
		WillReturnRows(sqlmock.NewRows(authorColumns)).
		RowsWillBeClosed()

	tx, _ := db.Begin()
	ctx := context.Background()

	bookRepo := repository.NewBookRepository()
	authorResult, err := bookRepo.FindAuthor(ctx, tx, inputAuthor1.AuthorId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Error(t, err)
	require.Empty(t, authorResult)
}

func TestRepoBookFindByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBook2.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[0].AuthorId, inputBook2.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[1].AuthorId, inputBook2.Authors[1].Name))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookResult, err := bookRepo.FindBookById(ctx, tx, inputBook2.BookId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
	require.Equal(t, inputBook2, bookResult)
}

func TestRepoBookFindByIdFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBook2.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns)).RowsWillBeClosed()

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookResult, err := bookRepo.FindBookById(ctx, tx, inputBook2.BookId)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Error(t, err)
	require.Empty(t, bookResult)
}

func TestRepoBookFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook1.BookId, inputBook1.Title, inputBook1.Edition, inputBook1.Authors[0].AuthorId, inputBook1.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[0].AuthorId, inputBook2.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[1].AuthorId, inputBook2.Authors[1].Name))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	booksResult := bookRepo.FindAllBooks(ctx, tx)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.Len(t, booksResult, 2)
	require.Equal(t, inputBook1, booksResult[0])
	require.Equal(t, inputBook2, booksResult[1])
}
