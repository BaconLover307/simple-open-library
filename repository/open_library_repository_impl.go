package repository

import (
	"context"
	"simple-open-library/exception"
	"simple-open-library/lib"
	"simple-open-library/model/domain"
)

type OpenLibraryRepositoryImpl struct {
	OpenLibraryLib lib.OpenLibraryLib
}

func NewOpenLibraryRepositoryImpl(openLibrary lib.OpenLibraryLib) LibraryRepository {
	return &OpenLibraryRepositoryImpl{
		OpenLibraryLib: openLibrary,
	}
}

func (openLibRepo OpenLibraryRepositoryImpl) Subjects(ctx context.Context, subject string, page int) ([]domain.Book, error) {
	openLibSubjectsResponse := openLibRepo.OpenLibraryLib.BrowseSubjects(ctx, subject, page)
	
	var libraryBooks []domain.Book
	if openLibSubjectsResponse.WorkCount == 0 {
		return libraryBooks, exception.NewNotFoundError("subject not found")
	}
	for _, book := range openLibSubjectsResponse.Works {
		libraryBooks = append(libraryBooks, domain.NewBookFromOpenLibrary(&book))
	}
	return libraryBooks, nil
}