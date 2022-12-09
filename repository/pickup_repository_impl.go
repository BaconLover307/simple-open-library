package repository

import (
	"context"
	"database/sql"
	"simple-open-library/model/domain"
)

type PickupRepositoryImpl struct {
}

func (repo PickupRepositoryImpl) Schedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	panic("not implemented")
}

func (repo PickupRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	panic("not implemented")
}

func (repo PickupRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	panic("not implemented")
}

func (repo PickupRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (domain.Pickup, error) {
	panic("not implemented")
}

func (repo PickupRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Pickup {
	panic("not implemented")
}

