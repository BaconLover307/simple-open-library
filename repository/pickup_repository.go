package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type PickupRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	UpdateSchedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup)
	FindById(ctx context.Context, tx *sql.Tx, pickupId int) (domain.Pickup, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Pickup
}