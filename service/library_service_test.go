package service_test

import (
	"simple-open-library/model/web"
	"testing"

	"github.com/stretchr/testify/require"
)
var (
	inputSubjectSuccess = "science"
	inputSubjectNotFound = "asdfasdf"
)

func TestServiceLibraryBrowseBySubject(t *testing.T) {

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