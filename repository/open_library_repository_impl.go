package repository

import (
	"context"
	"simple-open-library/exception"
	lib "simple-open-library/lib"
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

func (openLibRepo OpenLibraryRepositoryImpl) Subjects(ctx context.Context, subject string, page int) ([]domain.LibraryBook, error) {
	libraryBooksResponse := openLibRepo.OpenLibraryLib.BrowseSubjects(ctx, subject, page)
	
	var libraryBooks []domain.LibraryBook
	if libraryBooksResponse.WorkCount == 0 {
		return libraryBooks, exception.NewNotFoundError("subject not found")
	}
	for _, book := range libraryBooksResponse.Works {
		libraryBooks = append(libraryBooks, domain.LibraryBook(domain.NewLibraryBookFromOpenLibrary(&book)))
	}
	return libraryBooks, nil
}