package web

import (
	"simple-open-library/model/domain"
	"sort"
)

type AuthorRequest struct {
	AuthorId string `validate:"min=1,max=20" json:"authorId"`
	Name     string `validate:"min=1,max=100" json:"name"`
}

func NewAuthor(request *AuthorRequest) domain.Author {
	return domain.Author{
		AuthorId: request.AuthorId,
		Name:     request.Name}
}

func NewAuthors(requests []AuthorRequest) []domain.Author {
	var authors []domain.Author

	for _, request := range requests {
		authors = append(authors, NewAuthor(&request))
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].AuthorId < authors[j].AuthorId
	})
	return authors
}
