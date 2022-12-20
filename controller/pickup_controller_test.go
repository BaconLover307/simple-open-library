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
	"time"

	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/test"

	"github.com/stretchr/testify/require"
)

var (
	inputPickup1 = domain.Pickup{
		PickupId: 1,
		Book:     inputBook1,
		Schedule: time.Now().Round(time.Second),
	}
	inputPickup1x = domain.Pickup{
		PickupId: 1,
		Book:     inputBook1x,
		Schedule: time.Now().Round(time.Second),
	}
	inputPickup2 = domain.Pickup{
		PickupId: 2,
		Book:     inputBook2,
		Schedule: time.Now().Round(time.Second),
	}
	inputSchedule1 = web.PickupUpdateScheduleRequest{
		Schedule: time.Now().Round(time.Second),
	}
	inputPickupBadSchedule = struct {
		PickupId int         `json:"pickup_id"`
		Book     domain.Book `json:"book"`
		Schedule string      `json:"schedule"`
	}{
		PickupId: 1,
		Book:     inputBook1x,
		Schedule: "2022-12-30T08:30:00",
	}
)

func TestControllerPickupSubmitSuccess(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)
	router := test.InitializeTestRouter(db)

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
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusOK, responseBody.Code)
	require.Equal(t, "OK", responseBody.Status)

	// $ Test body data
	jsonString, _ := json.Marshal(responseBody.Data)
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponse)

	require.Equal(t, inputPickup1.PickupId, pickupResponse.PickupId)
	require.Equal(t, inputPickup1.Schedule, pickupResponse.Schedule)

	require.Equal(t, inputPickup1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, inputPickup1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, inputPickup1.Book.Edition, pickupResponse.Book.Edition)

	require.Equal(t, inputPickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupSubmitFailedConflict(t *testing.T) {

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

	router := test.InitializeTestRouter(db)

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
	require.Equal(t, http.StatusConflict, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusConflict, responseBody.Code)
	require.Equal(t, "CONFLICT", responseBody.Status)
	require.Equal(t, "cannot overwrite existing book. please insert correct book data", responseBody.Data)
}

func TestControllerPickupSubmitFailedBadRequest(t *testing.T) {

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

	router := test.InitializeTestRouter(db)

	data, err := json.Marshal(inputPickupBadSchedule)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/pickups", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusBadRequest, responseBody.Code)
	require.Equal(t, "BAD REQUEST", responseBody.Status)
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

	router := test.InitializeTestRouter(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

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
	var pickupResponses []web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponses)

	require.Len(t, pickupResponses, 2)

	search := func(responses []web.PickupResponse, id int) web.PickupResponse {
		for _, response := range responses {
			if response.PickupId == id {
				return response
			}
		}
		return web.NewPickupResponse(&inputPickup1)
	}

	pickupOutput1 := search(pickupResponses, pickup1.PickupId)
	pickupOutput2 := search(pickupResponses, pickup2.PickupId)

	require.Equal(t, pickup1.PickupId, pickupOutput1.PickupId)
	require.Equal(t, pickup1.Schedule, pickupOutput1.Schedule)

	require.Equal(t, pickup1.Book.BookId, pickupOutput1.Book.BookId)
	require.Equal(t, pickup1.Book.Title, pickupOutput1.Book.Title)
	require.Equal(t, pickup1.Book.Edition, pickupOutput1.Book.Edition)

	require.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupOutput1.Book.Authors[0].AuthorId)
	require.Equal(t, pickup1.Book.Authors[0].Name, pickupOutput1.Book.Authors[0].Name)

	require.Equal(t, pickup2.PickupId, pickupOutput2.PickupId)
	require.Equal(t, pickup2.Schedule, pickupOutput2.Schedule)
	require.Equal(t, pickup2.Book.BookId, pickupOutput2.Book.BookId)
	require.Equal(t, pickup2.Book.Title, pickupOutput2.Book.Title)
	require.Equal(t, pickup2.Book.Edition, pickupOutput2.Book.Edition)
	require.Equal(t, pickup2.Book.Authors[0].AuthorId, pickupOutput2.Book.Authors[0].AuthorId)
	require.Equal(t, pickup2.Book.Authors[0].Name, pickupOutput2.Book.Authors[0].Name)
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

	router := test.InitializeTestRouter(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups/"+strconv.Itoa(1), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

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
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponse)

	require.Equal(t, pickup1.PickupId, pickupResponse.PickupId)
	require.Equal(t, pickup1.Schedule, pickupResponse.Schedule)

	require.Equal(t, pickup1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, pickup1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, pickup1.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, pickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
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

	router := test.InitializeTestRouter(db)

	request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/pickups/"+strconv.Itoa(2), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusNotFound, responseBody.Code)
	require.Equal(t, "NOT FOUND", responseBody.Status)

	// $ Test body data
	require.Equal(t, "pick up schedule not found", responseBody.Data)
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

	router := test.InitializeTestRouter(db)

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
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusOK, responseBody.Code)
	require.Equal(t, "OK", responseBody.Status)

	// $ Test body data
	jsonString, _ := json.Marshal(responseBody.Data)
	var pickupResponse web.PickupResponse
	json.Unmarshal(jsonString, &pickupResponse)

	require.Equal(t, pickup1.PickupId, pickupResponse.PickupId)
	require.Equal(t, inputSchedule1.Schedule, pickupResponse.Schedule)

	require.Equal(t, pickup1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, pickup1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, pickup1.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, pickup1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, pickup1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestControllerPickupUpdateFailedNotFound(t *testing.T) {

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

	router := test.InitializeTestRouter(db)

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
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusNotFound, responseBody.Code)
	require.Equal(t, "NOT FOUND", responseBody.Status)

	// $ Test body data
	require.Equal(t, "pick up schedule not found", responseBody.Data)
}

func TestControllerPickupUpdateFailedBadRequest(t *testing.T) {

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

	router := test.InitializeTestRouter(db)

	data, err := json.Marshal(inputPickupBadSchedule)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPut, BaseURL+"/api/pickups/"+strconv.Itoa(2), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusBadRequest, responseBody.Code)
	require.Equal(t, "BAD REQUEST", responseBody.Status)
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

	router := test.InitializeTestRouter(db)

	request := httptest.NewRequest(http.MethodDelete, BaseURL+"/api/pickups/"+strconv.Itoa(1), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

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
	require.Nil(t, responseBody.Data)
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

	router := test.InitializeTestRouter(db)

	request := httptest.NewRequest(http.MethodDelete, BaseURL+"/api/pickups/"+strconv.Itoa(2), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// $ Test HTTP status
	response := recorder.Result()
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	// $ Test body status & code
	require.Equal(t, http.StatusNotFound, responseBody.Code)
	require.Equal(t, "NOT FOUND", responseBody.Status)

	// $ Test body data
	require.Equal(t, "pick up schedule not found", responseBody.Data)
}

func TestControllerPickupUnauthorized(t *testing.T) {

	db := test.SetupTestDB()
	test.TruncateDatabase(db)
	router := test.InitializeTestRouter(db)

	data, err := json.Marshal(inputPickup1)
	helper.PanicIfError(err)
	requestBody := bytes.NewReader(data)
	request := httptest.NewRequest(http.MethodPost, BaseURL+"/api/pickups", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SALAH")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody web.WebResponse
	json.Unmarshal(body, &responseBody)

	require.Equal(t, http.StatusUnauthorized, responseBody.Code)
	require.Equal(t, "UNAUTHORIZED", responseBody.Status)
	require.Equal(t, nil, responseBody.Data)
}
