package service_test

import (
	"simple-open-library/model/web"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
		inputAuthor1 = web.AuthorRequest{
			AuthorId: "ta001",
			Name: "Jk Rolling",
		}
		inputAuthor2 = web.AuthorRequest{
			AuthorId: "ta002",
			Name: "Mark Mansion",
		}	
		inputAuthors1 = []web.AuthorRequest{inputAuthor1}
		inputAuthors2 = []web.AuthorRequest{inputAuthor1, inputAuthor2}
		
		inputBook1 = web.BookRequest{
			BookId: "tb001",
			Title: "Test Book",
			Edition: 1,
			Authors: inputAuthors1,
		}
		inputBook2 = web.BookRequest{
			BookId: "tb002",
			Title: "Help Book",
			Edition: 2,
			Authors: inputAuthors2,
		}
	)

func TestServiceBookSave(t *testing.T) {
	bookResponse := testBookService.SaveBook(testCtx, inputBook2)
	require.Equal(t, inputBook2.BookId, bookResponse.BookId)
	require.Equal(t, inputBook2.Title, bookResponse.Title)
	require.Equal(t, inputBook2.Edition, bookResponse.Edition)
	require.Equal(t, inputBook2.Authors[0].AuthorId, bookResponse.Authors[0].AuthorId)
	require.Equal(t, inputBook2.Authors[0].Name, bookResponse.Authors[0].Name)
	require.Equal(t, inputBook2.Authors[1].AuthorId, bookResponse.Authors[1].AuthorId)
	require.Equal(t, inputBook2.Authors[1].Name, bookResponse.Authors[1].Name)
}

func TestServiceBookFindById(t *testing.T) {
	bookResponse := testBookService.FindBookById(testCtx, inputBook2.BookId)
	require.Equal(t, inputBook2.BookId, bookResponse.BookId)
	require.Equal(t, inputBook2.Title, bookResponse.Title)
	require.Equal(t, inputBook2.Edition, bookResponse.Edition)
	require.Equal(t, inputBook2.Authors[0].AuthorId, bookResponse.Authors[0].AuthorId)
	require.Equal(t, inputBook2.Authors[0].Name, bookResponse.Authors[0].Name)
	require.Equal(t, inputBook2.Authors[1].AuthorId, bookResponse.Authors[1].AuthorId)
	require.Equal(t, inputBook2.Authors[1].Name, bookResponse.Authors[1].Name)
}

func TestServiceBookFindAll(t *testing.T) {
	testBookService.SaveBook(testCtx, inputBook1)
	bookResponses := testBookService.FindAllBooks(testCtx)

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
