package service_test

import (
	"context"
	"errors"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/repository"
	"simple-open-library/service"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

var (
	inputPickupCreate1 = web.PickupCreateRequest{
		Book:     inputBook1,
		Schedule: time.Now(),
	}
	inputPickupUpdate1 = web.PickupUpdateScheduleRequest{
		PickupId: 1,
		Schedule: time.Now(),
	}
	inputPickupCreate3 = web.PickupCreateRequest{
		Book:     inputBook2,
		Schedule: time.Now(),
	}
	inputPickupCreate4 = web.PickupCreateRequest{
		Book:     inputBook1,
		Schedule: time.Now(),
	}
	pickupColumns = []string{"pickupId", "schedule", "bookId", "title", "edition", "authorId", "name"}
)

func TestServicePickupCreate(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO pickup").
		WithArgs(inputPickupCreate1.Book.BookId, inputPickupCreate1.Schedule).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())
	pickupResponse := testPickupService.Create(ctx, inputPickupCreate1)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Equal(t, 1, pickupResponse.PickupId)
	require.Equal(t, inputPickupCreate1.Schedule, pickupResponse.Schedule)

	require.Equal(t, inputPickupCreate1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, inputPickupCreate1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, inputPickupCreate1.Book.Edition, pickupResponse.Book.Edition)
	require.Len(t, inputPickupCreate1.Book.Authors, 1)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestServicePickupUpdateSchedule(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(pickupColumns).
			AddRow(
				inputPickupUpdate1.PickupId,
				inputPickupCreate3.Schedule,
				inputPickupCreate3.Book.BookId,
				inputPickupCreate3.Book.Title,
				inputPickupCreate3.Book.Edition,
				inputPickupCreate3.Book.Authors[0].AuthorId,
				inputPickupCreate3.Book.Authors[0].Name,
			))
	mock.ExpectExec("UPDATE pickup SET schedule").
		WithArgs(inputPickupUpdate1.Schedule, inputPickupUpdate1.PickupId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())
	pickupResponse := testPickupService.UpdateSchedule(ctx, inputPickupUpdate1)

	require.Equal(t, inputPickupUpdate1.PickupId, pickupResponse.PickupId)
	require.Equal(t, inputPickupUpdate1.Schedule, pickupResponse.Schedule)

	require.Equal(t, inputPickupCreate3.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, inputPickupCreate3.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, inputPickupCreate3.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, inputPickupCreate3.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickupCreate3.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestServicePickupDeleteSuccess(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO pickup").
		WithArgs(inputPickupCreate1.Book.BookId, inputPickupCreate1.Schedule).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(pickupColumns).
			AddRow(
				1,
				inputPickupCreate1.Schedule,
				inputPickupCreate1.Book.BookId,
				inputPickupCreate1.Book.Title,
				inputPickupCreate1.Book.Edition,
				inputPickupCreate1.Book.Authors[0].AuthorId,
				inputPickupCreate1.Book.Authors[0].Name,
			))
	mock.ExpectExec("DELETE FROM pickup").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(1).
		WillReturnError(errors.New("pick up schedule not found"))
	mock.ExpectRollback()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())
	testPickupService.Create(ctx, inputPickupCreate1)

	require.NotPanics(t, func() { testPickupService.Delete(ctx, 1) })
	require.PanicsWithError(t, "pick up schedule not found", func() { testPickupService.FindById(ctx, 1) })

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestServicePickupDeleteFailed(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WithArgs(1).
		WillReturnError(errors.New("pick up schedule not found"))
	mock.ExpectRollback()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())
	require.PanicsWithError(t, "pick up schedule not found", func() { testPickupService.Delete(ctx, 1) })

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}
}

func TestServicePickupFindById(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WillReturnRows(sqlmock.NewRows(pickupColumns).
			AddRow(1,
				inputPickupCreate3.Schedule,
				inputPickupCreate3.Book.BookId,
				inputPickupCreate3.Book.Title,
				inputPickupCreate3.Book.Edition,
				inputPickupCreate3.Book.Authors[0].AuthorId,
				inputPickupCreate3.Book.Authors[0].Name,
			).
			AddRow(1,
				inputPickupCreate3.Schedule,
				inputPickupCreate3.Book.BookId,
				inputPickupCreate3.Book.Title,
				inputPickupCreate3.Book.Edition,
				inputPickupCreate3.Book.Authors[1].AuthorId,
				inputPickupCreate3.Book.Authors[1].Name,
			).
			AddRow(2,
				inputPickupCreate4.Schedule,
				inputPickupCreate4.Book.BookId,
				inputPickupCreate4.Book.Title,
				inputPickupCreate4.Book.Edition,
				inputPickupCreate4.Book.Authors[0].AuthorId,
				inputPickupCreate4.Book.Authors[0].Name,
			))
	mock.ExpectCommit()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())
	pickupResponse := testPickupService.FindById(ctx, 2)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	input := inputPickupCreate4

	require.Equal(t, 2, pickupResponse.PickupId)
	require.Equal(t, input.Schedule, pickupResponse.Schedule)

	require.Equal(t, input.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, input.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, input.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, input.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, input.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestServicePickupFindAll(t *testing.T) {
	testDB, mock, err := sqlmock.New()
	helper.FatalIfMockError(t, err)
	defer testDB.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM pickup").
		WillReturnRows(sqlmock.NewRows(pickupColumns).
			AddRow(1,
				inputPickupCreate3.Schedule,
				inputPickupCreate3.Book.BookId,
				inputPickupCreate3.Book.Title,
				inputPickupCreate3.Book.Edition,
				inputPickupCreate3.Book.Authors[0].AuthorId,
				inputPickupCreate3.Book.Authors[0].Name,
			).
			AddRow(1,
				inputPickupCreate3.Schedule,
				inputPickupCreate3.Book.BookId,
				inputPickupCreate3.Book.Title,
				inputPickupCreate3.Book.Edition,
				inputPickupCreate3.Book.Authors[1].AuthorId,
				inputPickupCreate3.Book.Authors[1].Name,
			).
			AddRow(2,
				inputPickupCreate4.Schedule,
				inputPickupCreate4.Book.BookId,
				inputPickupCreate4.Book.Title,
				inputPickupCreate4.Book.Edition,
				inputPickupCreate4.Book.Authors[0].AuthorId,
				inputPickupCreate4.Book.Authors[0].Name,
			))
	mock.ExpectCommit()

	ctx := context.Background()
	testPickupService := service.NewPickupService(repository.NewPickupRepository(), testDB, validator.New())

	pickupResponses := testPickupService.FindAll(ctx)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectation error: %s", err)
	}

	require.Len(t, pickupResponses, 2)
	book1 := pickupResponses[0]
	book2 := pickupResponses[1]
	input1 := inputPickupCreate3
	input2 := inputPickupCreate4

	require.Equal(t, 1, book1.PickupId)
	require.Equal(t, input1.Schedule, book1.Schedule)

	require.Equal(t, input1.Book.BookId, book1.Book.BookId)
	require.Equal(t, input1.Book.Title, book1.Book.Title)
	require.Equal(t, input1.Book.Edition, book1.Book.Edition)
	require.Equal(t, input1.Book.Authors[0].AuthorId, book1.Book.Authors[0].AuthorId)
	require.Equal(t, input1.Book.Authors[0].Name, book1.Book.Authors[0].Name)
	require.Equal(t, input1.Book.Authors[1].AuthorId, book1.Book.Authors[1].AuthorId)
	require.Equal(t, input1.Book.Authors[1].Name, book1.Book.Authors[1].Name)

	require.Equal(t, 2, book2.PickupId)
	require.Equal(t, input2.Schedule, book2.Schedule)

	require.Equal(t, input2.Book.BookId, book2.Book.BookId)
	require.Equal(t, input2.Book.Title, book2.Book.Title)
	require.Equal(t, input2.Book.Edition, book2.Book.Edition)
	require.Equal(t, input2.Book.Authors[0].AuthorId, book2.Book.Authors[0].AuthorId)
	require.Equal(t, input2.Book.Authors[0].Name, book2.Book.Authors[0].Name)
}
