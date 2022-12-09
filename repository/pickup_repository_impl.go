package repository

import (
	"context"
	"database/sql"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
)

type PickupRepositoryImpl struct {
}

func NewPickupRepository() PickupRepository {
	return &PickupRepositoryImpl{}
}

func (repo PickupRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	query := "INSERT INTO pickup(bookId, schedule) VALUES(?, ?)"
	result, err := tx.ExecContext(ctx, query, pickup.Book.BookId, pickup.Schedule)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	pickup.PickupId = int(id)

	return pickup
}

func (repo PickupRepositoryImpl) UpdateSchedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	query := "UPDATE pickup set schedule = ? where pickupId = ?"
	_, err := tx.ExecContext(ctx, query, pickup.Schedule, pickup.PickupId)
	helper.PanicIfError(err)

	return pickup
}

func (repo PickupRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) {
	query := "DELETE FROM pickup WHERE pickupId = ?"
	_, err := tx.ExecContext(ctx, query, pickup.PickupId)
	helper.PanicIfError(err)

}

func (repo PickupRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, pickupId int) (domain.Pickup, error) {
	query := "SELECT p.pickupId, p.bookId, b.title, b.author, b.edition, p.schedule FROM pickup p JOIN book b ON p.bookId = b.bookId WHERE pickupId = ?"
	rows, err := tx.QueryContext(ctx, query, pickupId)
	helper.PanicIfError(err)
	defer rows.Close()

	pickup := domain.Pickup{}
	book := domain.Book{}
	if rows.Next() {
		err := rows.Scan(&pickup.PickupId, &book.BookId, &book.Title, &book.Author, &book.Edition, &pickup.Schedule)
		pickup.Book = book
		helper.PanicIfError(err)
		return pickup, nil
	} else {
		return pickup, exception.NewNotFoundError("pick up schedule found not found")
	}
}

func (repo PickupRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx, ) []domain.Pickup {
	query := "SELECT p.pickupId, p.bookId, b.title, b.author, b.edition, p.schedule FROM pickup p JOIN book b ON p.bookId = b.bookId"
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var pickups []domain.Pickup
	for rows.Next() {
		pickup := domain.Pickup{}
		book := domain.Book{}
		
		err := rows.Scan(&pickup.PickupId, &book.BookId, &book.Title, &book.Author, &book.Edition, &pickup.Schedule)
		pickup.Book = book
		helper.PanicIfError(err)

		pickups = append(pickups, pickup)
	}
	return pickups
}

