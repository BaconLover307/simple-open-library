package controller_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/test"
	"testing"

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
	inputBook1x = domain.Book{
		BookId:  "tb001",
		Title:   "Test Boo",
		Edition: 1,
		Authors: inputAuthors1,
	}
	inputBook2 = domain.Book{
		BookId:  "tb002",
		Title:   "Help Book",
		Edition: 2,
		Authors: inputAuthors2,
	}
)

func TestControllerBookListBooks(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {

		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestServer(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/books", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// $ Test HTTP status
		response := recorder.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)

		body, _ := io.ReadAll(response.Body)
		var responseBody web.WebResponse
		json.Unmarshal(body, &responseBody)

		// $ Test body status & code
		require.Equal(t, http.StatusOK, responseBody.Code)
		require.Equal(t, "OK", responseBody.Status)

		// $ Test body data
		jsonString, _ := json.Marshal(responseBody.Data)
		var bookResponses []web.BookResponse
		json.Unmarshal(jsonString, &bookResponses)

		require.Len(t, bookResponses, 0)
	})
	t.Run("Available", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)

		tx, _ := db.Begin()
		ctx := context.Background()
		bookRepo := repository.NewBookRepository()
		bookRepo.SaveBook(ctx, tx, inputBook1)
		bookRepo.SaveBook(ctx, tx, inputBook2)
		bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
		bookRepo.SaveAuthor(ctx, tx, inputAuthor2)
		bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
		bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook2.BookId)
		bookRepo.Authored(ctx, tx, inputAuthor2.AuthorId, inputBook2.BookId)
		tx.Commit()

		router := test.InitializeTestServer(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/books", nil)
		request.Header.Add("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// $ Test HTTP status
		response := recorder.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)

		body, _ := io.ReadAll(response.Body)
		var responseBody web.WebResponse
		json.Unmarshal(body, &responseBody)

		// $ Test body status & code
		require.Equal(t, http.StatusOK, responseBody.Code)
		require.Equal(t, "OK", responseBody.Status)

		// $ Test body data
		jsonString, _ := json.Marshal(responseBody.Data)
		var bookResponses []web.BookResponse
		json.Unmarshal(jsonString, &bookResponses)

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
	})
}
