package domain

import "simple-open-library/lib/model"

type LibraryBook struct {
	Title   string
	Author  string
	Edition int
}

func NewLibraryBookFromOpenLibrary(libraryBook *model.OpenLibraryBook) LibraryBook {
	return LibraryBook{
		Title: libraryBook.Title,
		Author: libraryBook.Authors[0].Name,
		Edition: libraryBook.EditionCount,
	}
}