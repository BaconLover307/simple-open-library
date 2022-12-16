package service_test

import (
	"context"
	"simple-open-library/helper"
	"simple-open-library/lib"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)
var (
	inputSubjectSuccess = "science"
	inputSubjectNotFound = "asdfasdf"
)

func TestServiceLibraryBrowseBySubject(t *testing.T) {
	// TODO
	testDB, _, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	testLibraryService := service.NewOpenLibraryService(lib.NewOpenLibraryLib(), testDB, validator.New())
	testCtx := context.Background()

	t.Run("Success", func(t *testing.T) {
		subjectResponse := testLibraryService.BrowseBySubject(testCtx, web.SubjectRequest{Subject: inputSubjectSuccess, Page: 1})
	
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 1, subjectResponse.Page)
	})

	t.Run("SuccessPage3", func(t *testing.T) {
		subjectResponse := testLibraryService.BrowseBySubject(testCtx, web.SubjectRequest{Subject: inputSubjectSuccess, Page: 3})
		
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 3, subjectResponse.Page)
	})

	t.Run("SubjectNotFound", func(t *testing.T) {
		require.PanicsWithError(t, "subject not found", func() {testLibraryService.BrowseBySubject(testCtx, web.SubjectRequest{Subject: inputSubjectNotFound, Page: 1})})
	})
}