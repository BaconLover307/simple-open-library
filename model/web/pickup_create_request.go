package web

type PickupCreateRequest struct {
	Book BookResponse	`validate:"required" json:"book"`
	Schedule string		`validate:"required" json:"schedule"`
}
