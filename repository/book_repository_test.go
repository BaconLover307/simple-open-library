package repository_test

import (
	"context"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
		inputAuthor1 = domain.Author{
			AuthorId: "ta001",
			Name: "Jk Rolling",
		}
		inputAuthor2 = domain.Author{
			AuthorId: "ta002",
			Name: "Mark Mansion",
		}	
		inputAuthors1 = []domain.Author{inputAuthor1}
		inputAuthors2 = []domain.Author{inputAuthor1, inputAuthor2}
		
		inputBook1 = domain.Book{
			BookId: "tb001",
			Title: "Test Book",
			Edition: 1,
			Authors: inputAuthors1,
		}
		inputBook2 = domain.Book{
			BookId: "tb002",
			Title: "Help Book",
			Edition: 2,
			Authors: inputAuthors2,
		}
	)

func TestRepoBookFindAuthor(t *testing.T) {
	testDB.Begin()
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveAuthor(ctx, testTx, inputAuthor1)
	authorResult, err := bookRepo.FindAuthor(ctx, testTx, inputAuthor1.AuthorId)
	
	require.NoError(t, err)
	require.Equal(t, inputAuthor1, authorResult)
	
	authorResult, err = bookRepo.FindAuthor(ctx, testTx, inputAuthor2.AuthorId)
	require.Error(t, err)
	require.Empty(t, authorResult)
}

func TestRepoBookFindBookById(t *testing.T) {
	testDB.Begin()
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)
	
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, testTx, inputBook1)
	bookRepo.Authored(ctx, testTx, inputBook1.Authors[0].AuthorId, inputBook1.BookId)
	bookResult, err := bookRepo.FindBookById(ctx, testTx, inputBook1.BookId)
	
	require.NoError(t, err)
	require.Equal(t, inputBook1, bookResult)
	
	bookResult, err = bookRepo.FindBookById(ctx, testTx, inputBook2.BookId)
	
	require.Error(t, err)
	require.Empty(t, bookResult)
}


func TestRepoBookFindAll(t *testing.T) {
	testDB.Begin()
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)

	defer helper.CommitOrRollback(testTx)
	

	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveAuthor(ctx, testTx, inputAuthor2)
	bookRepo.SaveBook(ctx, testTx, inputBook2)
	bookRepo.Authored(ctx, testTx, inputAuthor1.AuthorId, inputBook2.BookId)
	bookRepo.Authored(ctx, testTx, inputAuthor2.AuthorId, inputBook2.BookId)

	booksResult := bookRepo.FindAllBooks(ctx, testTx)

	require.Equal(t, 2, len(booksResult))
	require.Equal(t, inputBook1, booksResult[0])
	require.Equal(t, inputBook2, booksResult[1])

}