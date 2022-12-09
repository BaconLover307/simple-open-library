package service

import (
	"context"
	"simple-open-library/model/web"
)

type PickupService interface {
	Create(ctx context.Context, request web.PickupCreateRequest) web.PickupResponse
	UpdateSchedule(ctx context.Context, request web.PickupUpdateScheduleRequest) web.PickupResponse
	Delete(ctx context.Context, pickupId int)
	FindById(ctx context.Context, pickupId int) web.PickupResponse
	FindAll(ctx context.Context) []web.PickupResponse
}