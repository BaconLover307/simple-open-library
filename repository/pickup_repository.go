package repository

import (
	"context"
	"database/sql"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"time"
)

type PickupRepository interface {
	Create(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	UpdateSchedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup
	Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup)
	FindById(ctx context.Context, tx *sql.Tx, pickupId int) (domain.Pickup, error)
	FindAll(ctx context.Context, tx *sql.Tx) []domain.Pickup
}

type pickupRepo struct {
}

func NewPickupRepository() PickupRepository {
	return &pickupRepo{}
}

func (repo pickupRepo) Create(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	query := "INSERT INTO pickup(bookId, schedule) VALUES(?, ?)"
	result, err := tx.ExecContext(ctx, query, pickup.Book.BookId, pickup.Schedule.Round(time.Second))
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	pickup.PickupId = int(id)

	return pickup
}

func (repo pickupRepo) UpdateSchedule(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) domain.Pickup {
	query := "UPDATE pickup SET schedule = ? where pickupId = ?"
	_, err := tx.ExecContext(ctx, query, pickup.Schedule.Round(time.Second), pickup.PickupId)
	helper.PanicIfError(err)

	return pickup
}

func (repo pickupRepo) Delete(ctx context.Context, tx *sql.Tx, pickup domain.Pickup) {
	query := "DELETE FROM pickup WHERE pickupId = ?"
	_, err := tx.ExecContext(ctx, query, pickup.PickupId)
	helper.PanicIfError(err)

}

func (repo pickupRepo) FindById(ctx context.Context, tx *sql.Tx, pickupId int) (domain.Pickup, error) {
	query := `
	SELECT p.pickupId, p.schedule, p.bookId, b.title, b.edition, a.authorId, a.name FROM pickup p JOIN book b ON p.bookId = b.bookId
		LEFT JOIN authored ab ON b.bookId = ab.bookId
		LEFT JOIN author a ON ab.authorId = a.authorId
	WHERE p.pickupId = ?
	`
	rows, err := tx.QueryContext(ctx, query, pickupId)
	helper.PanicIfError(err)
	defer rows.Close()

	pickup := domain.Pickup{}
	book := domain.Book{}
	if rows.Next() {
		var authors []domain.Author
		author := domain.Author{}
		err = rows.Scan(&pickup.PickupId, &pickup.Schedule, &book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
		helper.PanicIfError(err)
		authors = append(authors, author)

		// Get remaining authors
		for rows.Next() {
			author := domain.Author{}
			err = rows.Scan(&pickup.PickupId, &pickup.Schedule, &book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
			helper.PanicIfError(err)
			authors = append(authors, author)
		}
		book.Authors = authors
		pickup.Book = book
		pickup.Schedule = pickup.Schedule.Local()
		helper.PanicIfError(err)

		return pickup, nil
	} else {
		return pickup, exception.NewNotFoundError("pick up schedule not found")
	}
}

func (repo pickupRepo) FindAll(ctx context.Context, tx *sql.Tx) []domain.Pickup {
	query := `
	SELECT p.pickupId, p.schedule, p.bookId, b.title, b.edition, a.authorId, a.name FROM pickup p JOIN book b ON p.bookId = b.bookId
		LEFT JOIN authored ab ON b.bookId = ab.bookId
		LEFT JOIN author a ON ab.authorId = a.authorId
		ORDER BY p.pickupId
	`
	rows, err := tx.QueryContext(ctx, query)
	helper.PanicIfError(err)
	defer rows.Close()

	var pickups []domain.Pickup
	pickup := domain.Pickup{}
	book := domain.Book{}
	author := domain.Author{}
	var authors []domain.Author

	for rows.Next() {
		err = rows.Scan(&pickup.PickupId, &pickup.Schedule, &book.BookId, &book.Title, &book.Edition, &author.AuthorId, &author.Name)
		helper.PanicIfError(err)
		if len(pickups) == 0 || pickups[len(pickups)-1].PickupId != pickup.PickupId {
			authors = nil
			authors = append(authors, author)
			book.Authors = authors
			pickup.Book = book
			pickup.Schedule = pickup.Schedule.Local()
			helper.PanicIfError(err)
			pickups = append(pickups, pickup)
		} else {
			authors = append(authors, author)
			book.Authors = authors
			pickup.Book = book
			pickup.Schedule = pickup.Schedule.Local()
			helper.PanicIfError(err)
			pickups[len(pickups)-1] = pickup
		}
	}
	return pickups
}
