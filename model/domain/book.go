package domain

import (
	"simple-open-library/lib/model"
	"strings"
)

type Book struct {
	BookId  string
	Title   string
	Authors []Author
	Edition int
}

const (
	openLibraryAuthorKeyPrefix = "/authors/"
	openLibraryBookKeyPrefix   = "/works/"
)

func NewBookFromOpenLibrary(libraryBook *model.OpenLibraryBook) Book {
	var authors []Author
	for _, author := range libraryBook.Authors {
		authorId := strings.TrimPrefix(author.Key, openLibraryAuthorKeyPrefix)
		authors = append(authors, Author{authorId, author.Name})
	}

	bookId := strings.TrimPrefix(libraryBook.Key, openLibraryBookKeyPrefix)

	return Book{
		BookId:  bookId,
		Title:   libraryBook.Title,
		Authors: authors,
		Edition: libraryBook.EditionCount,
	}
}
