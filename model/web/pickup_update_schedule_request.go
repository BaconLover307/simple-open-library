package web

import "time"

type PickupUpdateScheduleRequest struct {
	PickupId int       `validate:"required" json:"pickupId"`
	Schedule time.Time `validate:"required" json:"schedule"`
}
