package web

type PickupUpdateScheduleRequest struct {
	PickupId int	`validate:"required" json:"pickupId"`
	Schedule string	`validate:"required" json:"schedule"`
}
