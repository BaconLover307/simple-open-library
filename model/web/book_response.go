package web

import "simple-open-library/model/domain"

type BookResponse struct {
	BookId  int		`json:"bookId"`
	Title   string	`json:"title"`
	Author  string 	`json:"author"`
	Edition int		`json:"edition"`
}

func NewBookResponse(book *domain.Book) BookResponse {
	return BookResponse{
		BookId: book.BookId,
		Title: book.Title,
		Author: book.Author,
		Edition: book.Edition,
	}
}
