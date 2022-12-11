package service_test

import (
	"context"
	"simple-open-library/lib"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)
var (
	inputSubjectSuccess = "science"
	inputSubjectNotFound = "asdfasdf"
)

func TestBrowseBySubject(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		service := service.NewLibraryService(lib.NewOpenLibraryLib(), testDB, validator.New())
		subjectResponse := service.BrowseBySubject(context.Background(), web.SubjectRequest{Subject: inputSubjectSuccess, Page: 1})
	
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 1, subjectResponse.Page)
	})

	t.Run("SuccessPage3", func(t *testing.T) {
		service := service.NewLibraryService(lib.NewOpenLibraryLib(), testDB, validator.New())
		subjectResponse := service.BrowseBySubject(context.Background(), web.SubjectRequest{Subject: inputSubjectSuccess, Page: 3})
		
		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 3, subjectResponse.Page)
	})

	t.Run("SubjectNotFound", func(t *testing.T) {
		service := service.NewLibraryService(lib.NewOpenLibraryLib(), testDB, validator.New())
		assert.PanicsWithError(t, "subject not found", func() {service.BrowseBySubject(context.Background(), web.SubjectRequest{Subject: inputSubjectNotFound, Page: 1})})
	})
}