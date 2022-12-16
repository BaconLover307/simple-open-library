package web

import (
	"simple-open-library/model/domain"
	"time"
)

type PickupResponse struct {
	PickupId   int		`json:"pickupId"`
	Book BookResponse	`json:"book"`
	Schedule time.Time		`json:"schedule"`
}

func NewPickupResponse(pickup *domain.Pickup) PickupResponse {
	return PickupResponse{
		PickupId: pickup.PickupId,
		Book: NewBookResponse(&pickup.Book),
		Schedule: pickup.Schedule,
	}
}

func NewPickupResponses(pickups []domain.Pickup) []PickupResponse {
	var responses []PickupResponse
	for _, pickup := range pickups {
		responses = append(responses, NewPickupResponse(&pickup))
	}
	return responses
}

