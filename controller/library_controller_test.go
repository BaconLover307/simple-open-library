package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"simple-open-library/model/web"
	"simple-open-library/test"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inputSubjectSuccess = "science"
	inputSubjectSuccessPage2 = "science?page=2"
	inputSubjectNotFound = "asdfasdf"
)

func TestBrowseBySubject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestServer(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/" + inputSubjectSuccess, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// $ Test HTTP status
		response := recorder.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)
		
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		// $ Test body status & code
		require.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
		require.Equal(t, "OK", responseBody["status"])
	
		// $ Test body data
		jsonString, _ := json.Marshal(responseBody["data"])
		var subjectResponse web.SubjectResponse
		json.Unmarshal(jsonString,&subjectResponse)	
		
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 1, subjectResponse.Page)
	})

	t.Run("SuccessPage2", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestServer(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/" + inputSubjectSuccessPage2, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// $ Test HTTP status
		response := recorder.Result()
		require.Equal(t, http.StatusOK, response.StatusCode)
		
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		// $ Test body status & code
		require.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
		require.Equal(t, "OK", responseBody["status"])
		
		// $ Test body data
		jsonString, _ := json.Marshal(responseBody["data"])
		var subjectResponse web.SubjectResponse
		json.Unmarshal(jsonString,&subjectResponse)
	
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 2, subjectResponse.Page)
	})

	t.Run("SubjectNotFound", func(t *testing.T) {
		db := test.SetupTestDB()
		test.TruncateDatabase(db)
		router := test.InitializeTestServer(db)

		request := httptest.NewRequest(http.MethodGet, BaseURL+"/api/subjects/" + inputSubjectNotFound, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("X-API-KEY", os.Getenv("X-API-KEY"))

		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)

		// $ Test HTTP status
		response := recorder.Result()
		require.Equal(t, http.StatusNotFound, response.StatusCode)
		
		body, _ := io.ReadAll(response.Body)
		var responseBody map[string]interface{}
		json.Unmarshal(body, &responseBody)

		// $ Test body status & code
		require.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
		require.Equal(t, "NOT FOUND", responseBody["status"])
	
		// $ Test body data
		subjectResponse := responseBody["data"].(string)
	
		require.Equal(t, "subject not found", subjectResponse)
	})
}
