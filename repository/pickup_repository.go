package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type PickupRepository interface {
	Schedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	Update(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pickup, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Pickup
}