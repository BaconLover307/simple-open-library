package repository_test

import (
	"context"
	"errors"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

var (
	inputPickup1 = domain.Pickup{
		PickupId: 1,
		Book:     inputBook1,
		Schedule: time.Now().Round(time.Second),
	}
	inputPickup2 = domain.Pickup{
		PickupId: 2,
		Book:     inputBook2,
		Schedule: time.Now().Round(time.Second),
	}
	inputPickup3 = domain.Pickup{
		PickupId: 3,
		Book:     inputBook2,
		Schedule: time.Now().Round(time.Second),
	}
	inputSchedule1 = time.Now().Round(time.Second).Add(1000*time.Second)

	selectPickupColumns = []string{"pickupId", "schedule", "bookId", "title", "edition", "authorId", "name"}
)

func TestRepoPickupCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO pickup").
		WithArgs(inputPickup1.Book.BookId, inputPickup1.Schedule).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	savedPickup := pickupRepo.Create(ctx, tx, inputPickup1)
	inputPickup1.PickupId = savedPickup.PickupId

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	require.NoError(t, err)
	require.Equal(t, inputPickup1, savedPickup)
}

func TestRepoPickupFindByIdSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(selectPickupColumns).
			AddRow(
				inputPickup2.PickupId,
				inputPickup2.Schedule,
				inputPickup2.Book.BookId,
				inputPickup2.Book.Title,
				inputPickup2.Book.Edition,
				inputPickup2.Book.Authors[0].AuthorId,
				inputPickup2.Book.Authors[0].Name,
			).
			AddRow(
				inputPickup2.PickupId,
				inputPickup2.Schedule,
				inputPickup2.Book.BookId,
				inputPickup2.Book.Title,
				inputPickup2.Book.Edition,
				inputPickup2.Book.Authors[1].AuthorId,
				inputPickup2.Book.Authors[1].Name,
			))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupResult, err := pickupRepo.FindById(ctx, tx, inputPickup2.PickupId)

	require.NoError(t, err)
	require.Equal(t, inputPickup2, pickupResult)
}
func TestRepoPickupFindByIdFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(selectPickupColumns)).
		RowsWillBeClosed()

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupResult, err := pickupRepo.FindById(ctx, tx, 2)

	require.ErrorContains(t, err, "pick up schedule not found")
	require.Empty(t, pickupResult)
}

func TestRepoPickupUpdateSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO pickup").
		WithArgs(inputPickup1.Book.BookId, inputPickup1.Schedule).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE pickup SET schedule").
		WithArgs(inputSchedule1, inputPickup1.PickupId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	savedPickup := pickupRepo.Create(ctx, tx, inputPickup1)
	savedPickup.Schedule = inputSchedule1
	pickupResult := pickupRepo.UpdateSchedule(ctx, tx, savedPickup)

	require.NoError(t, err)
	require.Equal(t, savedPickup.Book, pickupResult.Book)
	require.Equal(t, savedPickup.PickupId, pickupResult.PickupId)
	require.Equal(t, inputSchedule1, pickupResult.Schedule)
}

func TestRepoPickupUpdateFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE pickup SET schedule").
		WithArgs(inputSchedule1, 1).
		WillReturnError(errors.New(""))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()

	require.Panics(t, func() { pickupRepo.UpdateSchedule(ctx, tx, inputPickup2) })
}

func TestRepoPickupDeleteSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM pickup").
		WithArgs(2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT (.+) FROM pickup p").
		WithArgs(inputPickup2.PickupId).
		WillReturnRows(sqlmock.NewRows(selectPickupColumns)).
		RowsWillBeClosed()

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Delete(ctx, tx, inputPickup2)
	pickupResult, err := pickupRepo.FindById(ctx, tx, inputPickup2.PickupId)

	require.Error(t, err)
	require.Empty(t, pickupResult)
}

func TestRepoPickupDeleteFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM pickup").
		WithArgs(2).WillReturnError(errors.New(""))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()

	require.Panics(t, func() { pickupRepo.Delete(ctx, tx, inputPickup2) })
}

func TestRepoPickupFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup p").
		WillReturnRows(sqlmock.NewRows(selectPickupColumns).
			AddRow(
				inputPickup1.PickupId,
				inputPickup1.Schedule,
				inputPickup1.Book.BookId,
				inputPickup1.Book.Title,
				inputPickup1.Book.Edition,
				inputPickup1.Book.Authors[0].AuthorId,
				inputPickup1.Book.Authors[0].Name,
			).
			AddRow(
				inputPickup3.PickupId,
				inputPickup3.Schedule,
				inputPickup3.Book.BookId,
				inputPickup3.Book.Title,
				inputPickup3.Book.Edition,
				inputPickup3.Book.Authors[0].AuthorId,
				inputPickup3.Book.Authors[0].Name,
			).
			AddRow(
				inputPickup3.PickupId,
				inputPickup3.Schedule,
				inputPickup3.Book.BookId,
				inputPickup3.Book.Title,
				inputPickup3.Book.Edition,
				inputPickup3.Book.Authors[1].AuthorId,
				inputPickup3.Book.Authors[1].Name,
			))

	tx, err := db.Begin()
	helper.PanicIfError(err)
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()

	pickupsResults := pickupRepo.FindAll(ctx, tx)

	require.Len(t, pickupsResults, 2)
	require.Equal(t, inputPickup1, pickupsResults[0])
	require.Equal(t, inputPickup3, pickupsResults[1])
}
