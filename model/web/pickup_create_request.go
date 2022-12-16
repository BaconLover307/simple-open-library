package web

import "time"

type PickupCreateRequest struct {
	Book BookRequest	`validate:"required" json:"book"`
	Schedule time.Time		`validate:"required" json:"schedule"`
}
