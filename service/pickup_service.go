package service

import (
	"context"
	"database/sql"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/model/web"
	"simple-open-library/repository"

	"github.com/go-playground/validator/v10"
)

type PickupService interface {
	Create(ctx context.Context, request web.PickupCreateRequest) web.PickupResponse
	UpdateSchedule(ctx context.Context, request web.PickupUpdateScheduleRequest) web.PickupResponse
	Delete(ctx context.Context, pickupId int)
	FindById(ctx context.Context, pickupId int) web.PickupResponse
	FindAll(ctx context.Context) []web.PickupResponse
}

type pickupService struct {
	Repo repository.PickupRepository
	DB *sql.DB
	Validate *validator.Validate
}

func NewPickupService(repo repository.PickupRepository, db *sql.DB, validate *validator.Validate) PickupService {
	return &pickupService{repo, db, validate}
}

func (service pickupService) Create(ctx context.Context, request web.PickupCreateRequest) web.PickupResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickup := domain.Pickup{
		Book: web.NewBook(&request.Book),
		Schedule: request.Schedule,
	}
	pickup = service.Repo.Create(ctx, tx, pickup)

	return web.NewPickupResponse(&pickup)
}

func (service pickupService) UpdateSchedule(ctx context.Context, request web.PickupUpdateScheduleRequest) web.PickupResponse {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickup, err := service.Repo.FindById(ctx, tx, request.PickupId)
	helper.PanicIfError(err)

	pickup.Schedule = request.Schedule
	pickup = service.Repo.UpdateSchedule(ctx, tx, pickup)

	return web.NewPickupResponse(&pickup)
}

func (service pickupService) Delete(ctx context.Context, pickupId int) {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickup, err := service.Repo.FindById(ctx, tx, pickupId)
	helper.PanicIfError(err)

	service.Repo.Delete(ctx, tx, pickup)
}

func (service pickupService) FindById(ctx context.Context, pickupId int) web.PickupResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickup, err := service.Repo.FindById(ctx, tx, pickupId)
	helper.PanicIfError(err)

	return web.NewPickupResponse(&pickup)
}

func (service pickupService) FindAll(ctx context.Context) []web.PickupResponse {
	tx, err := service.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	pickups := service.Repo.FindAll(ctx, tx)

	return web.NewPickupResponses(pickups)
}
