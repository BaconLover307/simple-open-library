package service_test

import (
	"context"
	"simple-open-library/lib/model"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var (
	inputSubjectSuccess  = "science"
	inputSubjectNotFound = "asdfasdf"

	sampleBook = model.OpenLibraryBook{
		Key:          "sample1",
		Title:        "Sample Library Book",
		EditionCount: 1,
		Authors:      []model.OpenLibraryAuthor{{Key: "author1", Name: "John Doe"}, {Key: "author2", Name: "Jane Doe"}},
	}
)

type stubOpenLibraryLib struct {
	browseMockLibrary func() model.OpenLibrarySubjectsResponse
}

func (stub *stubOpenLibraryLib) BrowseSubjects(ctx context.Context, subject string, page int) model.OpenLibrarySubjectsResponse {
	return stub.browseMockLibrary()
}

func TestServiceLibraryBrowseBySubject(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		browseFunc := func() model.OpenLibrarySubjectsResponse {
			return model.OpenLibrarySubjectsResponse{
				Name:      inputSubjectSuccess,
				WorkCount: 1000,
				Works:     []model.OpenLibraryBook{sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook},
			}
		}

		testLibraryService := service.NewOpenLibraryService(&stubOpenLibraryLib{browseMockLibrary: browseFunc}, validator.New())
		ctx := context.Background()
		subjectResponse := testLibraryService.BrowseBySubject(ctx, web.SubjectRequest{Subject: inputSubjectSuccess, Page: 1})

		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 1, subjectResponse.Page)
	})

	t.Run("SuccessPage3", func(t *testing.T) {
		browseFunc := func() model.OpenLibrarySubjectsResponse {
			return model.OpenLibrarySubjectsResponse{
				Name:      inputSubjectSuccess,
				WorkCount: 1000,
				Works:     []model.OpenLibraryBook{sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook, sampleBook},
			}
		}

		testLibraryService := service.NewOpenLibraryService(&stubOpenLibraryLib{browseMockLibrary: browseFunc}, validator.New())
		ctx := context.Background()

		subjectResponse := testLibraryService.BrowseBySubject(ctx, web.SubjectRequest{Subject: inputSubjectSuccess, Page: 3})

		require.Equal(t, inputSubjectSuccess, subjectResponse.Subject)
		require.Len(t, subjectResponse.Books, 10)
		require.Equal(t, 3, subjectResponse.Page)
	})

	t.Run("SubjectNotFound", func(t *testing.T) {
		browseFunc := func() model.OpenLibrarySubjectsResponse {
			return model.OpenLibrarySubjectsResponse{
				Name:      inputSubjectNotFound,
				WorkCount: 0,
				Works:     []model.OpenLibraryBook{},
			}
		}

		testLibraryService := service.NewOpenLibraryService(&stubOpenLibraryLib{browseMockLibrary: browseFunc}, validator.New())
		ctx := context.Background()

		require.PanicsWithError(t, "subject not found", func() {
			testLibraryService.BrowseBySubject(ctx, web.SubjectRequest{Subject: inputSubjectNotFound, Page: 1})
		})
	})
}
