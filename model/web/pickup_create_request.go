package web

type PickupCreateRequest struct {
	Book BookRequest	`validate:"required" json:"book"`
	Schedule string		`validate:"required" json:"schedule"`
}
