package service_test

import (
	"context"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var (
	inputAuthor1 = web.AuthorRequest{
		AuthorId: "ta001",
		Name:     "Jk Rolling",
	}
	inputAuthor2 = web.AuthorRequest{
		AuthorId: "ta002",
		Name:     "Mark Mansion",
	}
	inputAuthors1 = []web.AuthorRequest{inputAuthor1}
	inputAuthors2 = []web.AuthorRequest{inputAuthor1, inputAuthor2}

	inputBook1 = web.BookRequest{
		BookId:  "tb001",
		Title:   "Test Book",
		Edition: 1,
		Authors: inputAuthors1,
	}
	inputBook2 = web.BookRequest{
		BookId:  "tb002",
		Title:   "Help Book",
		Edition: 2,
		Authors: inputAuthors2,
	}
	inputBookOverwrite = web.BookRequest{
		BookId:  "tb001",
		Title:   "Test Bo",
		Edition: 1,
		Authors: inputAuthors1,
	}
	bookColumns       = []string{"bookId", "title", "edition"}
	authorColumns     = []string{"authorId", "name"}
	authoredColumns   = []string{"authorId", "bookId"}
	selectBookColumns = []string{"bookId", "title", "edition", "authorId", "name"}
)

func TestServiceBookSaveSuccess(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBook2.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns))
	mock.ExpectExec("INSERT INTO book").
		WithArgs(inputBook2.BookId, inputBook2.Title, inputBook2.Edition).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM author").
		WithArgs(inputBook2.Authors[0].AuthorId).
		WillReturnRows(sqlmock.NewRows(authorColumns))
	mock.ExpectExec("INSERT INTO author").
		WithArgs(inputBook2.Authors[0].AuthorId, inputBook2.Authors[0].Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO authored").
		WithArgs(inputBook2.Authors[0].AuthorId, inputBook2.BookId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM author").
		WithArgs(inputBook2.Authors[1].AuthorId).
		WillReturnRows(sqlmock.NewRows(authorColumns))
	mock.ExpectExec("INSERT INTO author").
		WithArgs(inputBook2.Authors[1].AuthorId, inputBook2.Authors[1].Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO authored").
		WithArgs(inputBook2.Authors[1].AuthorId, inputBook2.BookId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()
	validate := validator.New()
	bookService := service.NewBookService(repository.NewBookRepository(), testDB, validate)
	bookResponse := bookService.SaveBook(ctx, inputBook2)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Equal(t, inputBook2.BookId, bookResponse.BookId)
	require.Equal(t, inputBook2.Title, bookResponse.Title)
	require.Equal(t, inputBook2.Edition, bookResponse.Edition)
	require.Equal(t, inputBook2.Authors[0].AuthorId, bookResponse.Authors[0].AuthorId)
	require.Equal(t, inputBook2.Authors[0].Name, bookResponse.Authors[0].Name)
	require.Equal(t, inputBook2.Authors[1].AuthorId, bookResponse.Authors[1].AuthorId)
	require.Equal(t, inputBook2.Authors[1].Name, bookResponse.Authors[1].Name)
}

func TestServiceBookSaveSuccessExists(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBook1.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns))
	mock.ExpectExec("INSERT INTO book").
		WithArgs(inputBook1.BookId, inputBook1.Title, inputBook1.Edition).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM author").
		WithArgs(inputBook1.Authors[0].AuthorId).
		WillReturnRows(sqlmock.NewRows(authorColumns))
	mock.ExpectExec("INSERT INTO author").
		WithArgs(inputBook1.Authors[0].AuthorId, inputBook1.Authors[0].Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO authored").
		WithArgs(inputBook1.Authors[0].AuthorId, inputBook1.BookId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBook1.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook1.BookId, inputBook1.Title, inputBook1.Edition, inputBook1.Authors[0].AuthorId, inputBook1.Authors[0].Name))
	mock.ExpectCommit()

	ctx := context.Background()
	validate := validator.New()
	bookService := service.NewBookService(repository.NewBookRepository(), testDB, validate)

	bookResponse1 := bookService.SaveBook(ctx, inputBook1)
	bookResponse2 := bookService.SaveBook(ctx, inputBook1)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Equal(t, bookResponse1.BookId, bookResponse2.BookId)
	require.Equal(t, bookResponse1.Title, bookResponse2.Title)
	require.Equal(t, bookResponse1.Edition, bookResponse2.Edition)
	require.Equal(t, bookResponse1.Authors[0].AuthorId, bookResponse2.Authors[0].AuthorId)
	require.Equal(t, bookResponse1.Authors[0].Name, bookResponse2.Authors[0].Name)
}

func TestServiceBookSaveFailedOverwrite(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WithArgs(inputBookOverwrite.BookId).
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook1.BookId, inputBook1.Title, inputBook1.Edition, inputBook1.Authors[0].AuthorId, inputBook1.Authors[0].Name)).
		RowsWillBeClosed()
	mock.ExpectRollback()

	ctx := context.Background()
	validate := validator.New()
	bookService := service.NewBookService(repository.NewBookRepository(), testDB, validate)

	require.PanicsWithError(t, "cannot overwrite existing book. please insert correct book data", func() { bookService.SaveBook(ctx, inputBookOverwrite) })

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestServiceBookFindById(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[0].AuthorId, inputBook2.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[1].AuthorId, inputBook2.Authors[1].Name))
	mock.ExpectCommit()

	ctx := context.Background()
	validate := validator.New()
	testBookService := service.NewBookService(repository.NewBookRepository(), testDB, validate)
	bookResponse := testBookService.FindBookById(ctx, inputBook2.BookId)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Equal(t, inputBook2.BookId, bookResponse.BookId)
	require.Equal(t, inputBook2.Title, bookResponse.Title)
	require.Equal(t, inputBook2.Edition, bookResponse.Edition)
	require.Equal(t, inputBook2.Authors[0].AuthorId, bookResponse.Authors[0].AuthorId)
	require.Equal(t, inputBook2.Authors[0].Name, bookResponse.Authors[0].Name)
	require.Equal(t, inputBook2.Authors[1].AuthorId, bookResponse.Authors[1].AuthorId)
	require.Equal(t, inputBook2.Authors[1].Name, bookResponse.Authors[1].Name)
}

func TestServiceBookFindAll(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM author a JOIN authored ab").
		WillReturnRows(sqlmock.NewRows(selectBookColumns).
			AddRow(inputBook1.BookId, inputBook1.Title, inputBook1.Edition, inputBook1.Authors[0].AuthorId, inputBook1.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[0].AuthorId, inputBook2.Authors[0].Name).
			AddRow(inputBook2.BookId, inputBook2.Title, inputBook2.Edition, inputBook2.Authors[1].AuthorId, inputBook2.Authors[1].Name))
	mock.ExpectCommit()

	ctx := context.Background()
	validate := validator.New()
	testBookService := service.NewBookService(repository.NewBookRepository(), testDB, validate)
	bookResponses := testBookService.FindAllBooks(ctx)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Len(t, bookResponses, 2)
	book1 := bookResponses[1]
	require.Equal(t, inputBook2.BookId, book1.BookId)
	require.Equal(t, inputBook2.Title, book1.Title)
	require.Equal(t, inputBook2.Edition, book1.Edition)
	require.Equal(t, inputBook2.Authors[0].AuthorId, book1.Authors[0].AuthorId)
	require.Equal(t, inputBook2.Authors[0].Name, book1.Authors[0].Name)
	require.Equal(t, inputBook2.Authors[1].AuthorId, book1.Authors[1].AuthorId)
	require.Equal(t, inputBook2.Authors[1].Name, book1.Authors[1].Name)

	book2 := bookResponses[0]
	require.Equal(t, inputBook1.BookId, book2.BookId)
	require.Equal(t, inputBook1.Title, book2.Title)
	require.Equal(t, inputBook1.Edition, book2.Edition)
	require.Equal(t, inputBook1.Authors[0].AuthorId, book2.Authors[0].AuthorId)
	require.Equal(t, inputBook1.Authors[0].Name, book2.Authors[0].Name)

}
