package web

import "simple-open-library/model/domain"

type PickupResponse struct {
	PickupId   int		`json:"pickupId"`
	Book BookResponse	`json:"book"`
	Schedule string		`json:"schedule"`
}

func NewPickupResponse(pickup *domain.Pickup) PickupResponse {
	return PickupResponse{
		PickupId: pickup.PickupId,
		Book: BookResponse(pickup.Book),
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

