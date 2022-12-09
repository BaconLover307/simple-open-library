package web

type SubjectRequest struct {
	Subject string 	`validate:"min=1,max=200" json:"subject"`
	Page 	int 	`validate:"gte=1" json:"page"`
}