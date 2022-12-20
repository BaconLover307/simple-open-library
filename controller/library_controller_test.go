package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"simple-open-library/model/web"
	"simple-open-library/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	inputSubjectSuccess      = "science"
	inputSubjectSuccessPage2 = "science?page=2"
	inputSubjectNotFound     = "asdfasdf"
)

func TestControllerLibraryBrowseBySubject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestRouter(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/"+inputSubjectSuccess, nil)
		request.Header.Add("Content-Type", "application/json")

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
		var subjectResponse web.SubjectResponse
		json.Unmarshal(jsonString, &subjectResponse)

		assert.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		assert.Len(t, subjectResponse.Books, 10)
		assert.Equal(t, 1, subjectResponse.Page)
	})

	t.Run("SuccessPage2", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestRouter(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/"+inputSubjectSuccessPage2, nil)
		request.Header.Add("Content-Type", "application/json")

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
		var subjectResponse web.SubjectResponse
		json.Unmarshal(jsonString, &subjectResponse)

		assert.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		assert.Len(t, subjectResponse.Books, 10)
		assert.Equal(t, 2, subjectResponse.Page)
	})

	t.Run("SubjectNotFound", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestRouter(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/"+inputSubjectNotFound, nil)
		request.Header.Add("Content-Type", "application/json")

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
		subjectResponse := responseBody.Data.(string)

		assert.Equal(t, "subject not found", subjectResponse)
	})
}
