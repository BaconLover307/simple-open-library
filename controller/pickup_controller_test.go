package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/test"

	"github.com/stretchr/testify/assert"
)

var (
		inputAuthors = []domain.Author{
			{
				AuthorId: "ta001",
				Name: "Jk Rolling",
			},
		}
		
		inputBook = domain.Book{
			BookId: "tb001",
			Title: "Test Book",
			Edition: 1,
			Authors: inputAuthors,
		}

		inputPickup = domain.Pickup{
			PickupId: 1,
			Book: inputBook,
			Schedule: "2022-12-12 10:20:30",
		}
	)

func TestControllerPickupSubmit(t *testing.T) {

	data, err := json.Marshal(inputPickup)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/pickups", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	testRouter.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	
	// $ Test body data
	jsonString, _ := json.Marshal(responseBody.Data)
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponse)

	assert.Equal(t, inputPickup.PickupId, pickupResponse.PickupId)
	assert.Equal(t, inputPickup.Schedule, pickupResponse.Schedule)

	assert.Equal(t, inputPickup.Book.BookId, pickupResponse.Book.BookId)
	assert.Equal(t, inputPickup.Book.Title, pickupResponse.Book.Title)
	assert.Equal(t, inputPickup.Book.Edition, pickupResponse.Book.Edition)

	assert.Equal(t, inputPickup.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	assert.Equal(t, inputPickup.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupListSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook)
	bookRepo.SaveAuthor(ctx, tx, inputAuthors[0])
	bookRepo.Authored(ctx, tx, inputAuthors[0].AuthorId, inputBook.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickup1 := pickupRepo.Create(ctx, tx, inputPickup)
	pickup2 := pickupRepo.Create(ctx, tx, inputPickup)
	tx.Commit()

	router := test.InitializeTestServer(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusOK, responseBody.Code)
	assert.Equal(t, "OK", responseBody.Status)
	
	// $ Test body data
	jsonString, _ := json.Marshal(responseBody.Data)
	var pickupResponses []web.PickupResponse
	json.Unmarshal(jsonString,&pickupResponses)

	assert.Equal(t, pickup1.PickupId, pickupResponses[0].PickupId)
	assert.Equal(t, pickup1.Schedule, pickupResponses[0].Schedule)

	assert.Equal(t, pickup1.Book.BookId, pickupResponses[0].Book.BookId)
	assert.Equal(t, pickup1.Book.Title, pickupResponses[0].Book.Title)
	assert.Equal(t, pickup1.Book.Edition, pickupResponses[0].Book.Edition)

	assert.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupResponses[0].Book.Authors[0].AuthorId)
	assert.Equal(t, pickup1.Book.Authors[0].Name, pickupResponses[0].Book.Authors[0].Name)

	assert.Equal(t, pickup2.PickupId, pickupResponses[1].PickupId)
	assert.Equal(t, pickup2.Schedule, pickupResponses[1].Schedule)
	assert.Equal(t, pickup2.Book.BookId, pickupResponses[1].Book.BookId)
	assert.Equal(t, pickup2.Book.Title, pickupResponses[1].Book.Title)
	assert.Equal(t, pickup2.Book.Edition, pickupResponses[1].Book.Edition)
	assert.Equal(t, pickup2.Book.Authors[0].AuthorId, pickupResponses[1].Book.Authors[0].AuthorId)
	assert.Equal(t, pickup2.Book.Authors[0].Name, pickupResponses[1].Book.Authors[0].Name)
}

func TestControllerPickupUnauthorized(t *testing.T) {
	var (
		inputAuthors = []web.AuthorRequest{
			{
				AuthorId: "ta001",
				Name: "Jk Rolling",
			},
		}
		inputBook = web.BookRequest{
			BookId: "tb001",
			Title: "Test Book",
			Edition: 1,
			Authors: inputAuthors,
		}
		inputPickup = web.PickupCreateRequest{
			Book: inputBook,
			Schedule: "2022-12-12 10:20:30",
		}
	)
	
	db := test.SetupTestDB()
	test.TruncateDatabase(db)
	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputPickup)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SALAH")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusUnauthorized, responseBody.Code)
	assert.Equal(t, "UNAUTHORIZED", responseBody.Status)
	assert.Equal(t, nil, responseBody.Data)
}

