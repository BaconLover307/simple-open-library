package web

type SubjectResponse struct {
	Subject   string         `json:"subject"`
	BookCount int            `json:"bookCount"`
	Page      int            `json:"page"`
	Books     []BookResponse `json:"books"`
}