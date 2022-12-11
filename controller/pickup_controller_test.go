package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/test"

	"github.com/stretchr/testify/assert"
)

var (
		inputPickup1 = domain.Pickup{
			PickupId: 1,
			Book: inputBook1,
			Schedule: "2022-12-12 10:20:30",
		}
		inputPickup1x = domain.Pickup{
			PickupId: 1,
			Book: inputBook1x,
			Schedule: "2022-12-12 10:20:30",
		}
		inputPickup2 = domain.Pickup{
			PickupId: 2,
			Book: inputBook2,
			Schedule: "2022-12-12 10:20:30",
		}
		inputSchedule1 = web.PickupUpdateScheduleRequest{
			Schedule: "2001-02-02 02:22:20",
		}
	)

func TestControllerPickupSubmitSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)
	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputPickup1)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/pickups", requestBody)
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
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponse)

	assert.Equal(t, inputPickup1.PickupId, pickupResponse.PickupId)
	assert.Equal(t, inputPickup1.Schedule, pickupResponse.Schedule)

	assert.Equal(t, inputPickup1.Book.BookId, pickupResponse.Book.BookId)
	assert.Equal(t, inputPickup1.Book.Title, pickupResponse.Book.Title)
	assert.Equal(t, inputPickup1.Book.Edition, pickupResponse.Book.Edition)

	assert.Equal(t, inputPickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	assert.Equal(t, inputPickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupSubmitFailed(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputPickup1x)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/pickups", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusConflict, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusConflict, responseBody.Code)
	assert.Equal(t, "CONFLICT", responseBody.Status)
	assert.Equal(t, "cannot overwrite existing book. please insert correct book data", responseBody.Data)
}

func TestControllerPickupListSuccess(t *testing.T) {

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
	pickupRepo := repository.NewPickupRepository()
	pickup1 := pickupRepo.Create(ctx, tx, inputPickup1)
	pickup2 := pickupRepo.Create(ctx, tx, inputPickup2)
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

func TestControllerPickupGetSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickup1 := pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups/"+strconv.Itoa(1), nil)
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
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString,&pickupResponse)

	assert.Equal(t, pickup1.PickupId, pickupResponse.PickupId)
	assert.Equal(t, pickup1.Schedule, pickupResponse.Schedule)

	assert.Equal(t, pickup1.Book.BookId, pickupResponse.Book.BookId)
	assert.Equal(t, pickup1.Book.Title, pickupResponse.Book.Title)
	assert.Equal(t, pickup1.Book.Edition, pickupResponse.Book.Edition)
	assert.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	assert.Equal(t, pickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupGetFailed(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups/"+strconv.Itoa(2), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.Equal(t, "NOT FOUND", responseBody.Status)
	
	// $ Test body data
	assert.Equal(t, "pick up schedule not found", responseBody.Data)
}

func TestControllerPickupUpdateSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickup1 := pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputSchedule1)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPut, BaseURL+"/api/pickups/"+strconv.Itoa(1), requestBody)
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
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString,&pickupResponse)

	assert.Equal(t, pickup1.PickupId, pickupResponse.PickupId)
	assert.Equal(t, inputSchedule1.Schedule, pickupResponse.Schedule)

	assert.Equal(t, pickup1.Book.BookId, pickupResponse.Book.BookId)
	assert.Equal(t, pickup1.Book.Title, pickupResponse.Book.Title)
	assert.Equal(t, pickup1.Book.Edition, pickupResponse.Book.Edition)
	assert.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	assert.Equal(t, pickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupUpdateFailed(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputSchedule1)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPut, BaseURL+"/api/pickups/"+strconv.Itoa(2), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.Equal(t, "NOT FOUND", responseBody.Status)
	
	// $ Test body data
	assert.Equal(t, "pick up schedule not found", responseBody.Data)
}

func TestControllerPickupDeleteSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	request := httptest.NewRequest(http.MethodDelete, BaseURL+"/api/pickups/"+strconv.Itoa(1), nil)
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
	assert.Nil(t, responseBody.Data)
}

func TestControllerPickupDeleteFailed(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)

	tx, _ := db.Begin()
	ctx := context.Background()
	bookRepo := repository.NewBookRepository()
	bookRepo.SaveBook(ctx, tx, inputBook1)
	bookRepo.SaveAuthor(ctx, tx, inputAuthor1)
	bookRepo.Authored(ctx, tx, inputAuthor1.AuthorId, inputBook1.BookId)
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Create(ctx, tx, inputPickup1)
	tx.Commit()

	router := test.InitializeTestServer(db)

	request := httptest.NewRequest(http.MethodDelete, BaseURL+"/api/pickups/"+strconv.Itoa(2), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	
	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)
	
	// $ Test body status & code
	assert.Equal(t, http.StatusNotFound, responseBody.Code)
	assert.Equal(t, "NOT FOUND", responseBody.Status)
	
	// $ Test body data
	assert.Equal(t, "pick up schedule not found", responseBody.Data)
}


func TestControllerPickupUnauthorized(t *testing.T) {
	
	db := test.SetupTestDB()
	test.TruncateDatabase(db)
	router := test.InitializeTestServer(db)

	data, err := json.Marshal(inputPickup1)
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

