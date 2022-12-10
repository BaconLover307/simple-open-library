package web

import "simple-open-library/model/domain"

type AuthorResponse struct {
	AuthorId 	string 	`json:"authorId"`
	Name		string	`json:"name"`
}

func NewAuthorResponse(author *domain.Author) AuthorResponse {
	return AuthorResponse{
		AuthorId: author.AuthorId,
		Name: author.Name,
	}
}

func NewAuthorResponses(authors []domain.Author) []AuthorResponse {
	var authorResponses []AuthorResponse
	for _, author := range authors {
		authorResponses = append(authorResponses,NewAuthorResponse(&author))
	}
	return authorResponses
}

type BookResponse struct {
	BookId  string				`json:"bookId"`
	Title   string				`json:"title"`
	Authors  []AuthorResponse	`json:"authors"`
	Edition int					`json:"edition"`
}

func NewBookResponse(book *domain.Book) BookResponse {
	return BookResponse{
		BookId: book.BookId,
		Title: book.Title,
		Authors: NewAuthorResponses(book.Authors),
		Edition: book.Edition,
	}
}

func NewBookResponses(books []domain.Book) []BookResponse {
	var responses []BookResponse
	for _, book := range books {
		responses = append(responses, NewBookResponse(&book))
	}
	return responses
}
