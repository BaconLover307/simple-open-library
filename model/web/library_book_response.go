package web

import (
	"simple-open-library/model/domain"
)

type LibraryBookResponse struct {
	Title   string	`json:"title"`
	Author  string 	`json:"author"`
	Edition int		`json:"edition"`
}

func NewLibraryBookResponse(libraryBook *domain.LibraryBook) LibraryBookResponse {
	return LibraryBookResponse{
		Title: libraryBook.Title,
		Author: libraryBook.Author,
		Edition: libraryBook.Edition,
	}
}

func NewLibraryBookResponses(libraryBooks []domain.LibraryBook) []LibraryBookResponse {
	var responses []LibraryBookResponse
	for _, libraryBook := range libraryBooks {
		responses = append(responses, NewLibraryBookResponse(&libraryBook))
	}
	return responses
}
