package web

import "simple-open-library/model/domain"

type BookRequest struct {
	BookId  string          `validate:"required,min=1,max=20" json:"bookId"`
	Title   string          `validate:"required,min=1,max=200" json:"title"`
	Authors []AuthorRequest `validate:"required" json:"authors"`
	Edition int             `validate:"required,gte=1" json:"edition"`
}

func NewBook(request *BookRequest) domain.Book {
	return domain.Book{
		BookId:  request.BookId,
		Title:   request.Title,
		Authors: NewAuthors(request.Authors),
		Edition: request.Edition,
	}
}
