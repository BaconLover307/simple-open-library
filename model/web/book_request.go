package web

type BookRequest struct {
	Title 	string 	`validate:"min=1,max=200" json:"title"`
	Author 	string 	`validate:"min=1,max=200" json:"author"`
	Edition	int 	`json:"edition"`
}