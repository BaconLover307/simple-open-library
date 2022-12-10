package web

type SubjectResponse struct {
	Subject   string         `json:"subject"`
	BookCount int            `json:"book_count"`
	Page      int            `json:"page"`
	Books     []BookResponse `json:"books"`
}