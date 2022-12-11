package service_test

import (
	"simple-open-library/model/web"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
		inputPickupCreate1 = web.PickupCreateRequest{
			Book: inputBook1,
			Schedule: "2021-01-01 10:10:10",
		}
		inputPickupUpdate1 = web.PickupUpdateScheduleRequest{
			PickupId: 1,
			Schedule: "2022-02-02 12:20:20",
		}
		inputPickupCreate3 = web.PickupCreateRequest{
			Book: inputBook2,
			Schedule: "2013-01-03 13:23:33",
		}
		inputPickupCreate4 = web.PickupCreateRequest{
			Book: inputBook1,
			Schedule: "2012-02-02 02:22:20",
		}
	)

func TestServicePickupCreate(t *testing.T) {
	pickupResponse := testPickupService.Create(testCtx, inputPickupCreate1)
	
	require.Equal(t, 1, pickupResponse.PickupId)
	require.Equal(t, inputPickupCreate1.Schedule, pickupResponse.Schedule)

	require.Equal(t, inputPickupCreate1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, inputPickupCreate1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, inputPickupCreate1.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestServicePickupUpdateSchedule(t *testing.T) {
	pickupResponse := testPickupService.UpdateSchedule(testCtx, inputPickupUpdate1)

	require.Equal(t, inputPickupUpdate1.PickupId, pickupResponse.PickupId)
	require.Equal(t, inputPickupUpdate1.Schedule, pickupResponse.Schedule)

	require.Equal(t, inputPickupCreate1.Book.BookId, pickupResponse.Book.BookId)
	require.Equal(t, inputPickupCreate1.Book.Title, pickupResponse.Book.Title)
	require.Equal(t, inputPickupCreate1.Book.Edition, pickupResponse.Book.Edition)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].AuthorId, pickupResponse.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].Name, pickupResponse.Book.Authors[0].Name)
}

func TestServicePickupDelete(t *testing.T) {
	pickupResponse1 := testPickupService.FindById(testCtx, 1)
	require.Equal(t, 1, pickupResponse1.PickupId)
	require.Equal(t, inputPickupUpdate1.Schedule, pickupResponse1.Schedule)

	require.Equal(t, inputPickupCreate1.Book.BookId, pickupResponse1.Book.BookId)
	require.Equal(t, inputPickupCreate1.Book.Title, pickupResponse1.Book.Title)
	require.Equal(t, inputPickupCreate1.Book.Edition, pickupResponse1.Book.Edition)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].AuthorId, pickupResponse1.Book.Authors[0].AuthorId)
	require.Equal(t, inputPickupCreate1.Book.Authors[0].Name, pickupResponse1.Book.Authors[0].Name)

	testPickupService.Delete(testCtx, 1)
	require.PanicsWithError(t, "pick up schedule not found", func(){testPickupService.FindById(testCtx, 1)}) 
}

func TestServicePickupFindAll(t *testing.T) {
	pickupResponse1 := testPickupService.Create(testCtx, inputPickupCreate3)
	pickupResponse2 := testPickupService.Create(testCtx, inputPickupCreate4)
	pickupResponses := testPickupService.FindAll(testCtx)

	require.Len(t, pickupResponses, 2)
	book1 := pickupResponses[0]
	book2 := pickupResponses[1]
	
	require.Equal(t, pickupResponse1.PickupId, book1.PickupId)
	require.Equal(t, pickupResponse1.Schedule, book1.Schedule)
	
	require.Equal(t, pickupResponse1.Book.BookId, book1.Book.BookId)
	require.Equal(t, pickupResponse1.Book.Title, book1.Book.Title)
	require.Equal(t, pickupResponse1.Book.Edition, book1.Book.Edition)
	require.Equal(t, pickupResponse1.Book.Authors[0].AuthorId, book1.Book.Authors[0].AuthorId)
	require.Equal(t, pickupResponse1.Book.Authors[0].Name, book1.Book.Authors[0].Name)
	require.Equal(t, pickupResponse1.Book.Authors[1].AuthorId, book1.Book.Authors[1].AuthorId)
	require.Equal(t, pickupResponse1.Book.Authors[1].Name, book1.Book.Authors[1].Name)
		
	
	require.Equal(t, pickupResponse2.PickupId, book2.PickupId)
	require.Equal(t, pickupResponse2.Schedule, book2.Schedule)

	require.Equal(t, pickupResponse2.Book.BookId, book2.Book.BookId)
	require.Equal(t, pickupResponse2.Book.Title, book2.Book.Title)
	require.Equal(t, pickupResponse2.Book.Edition, book2.Book.Edition)
	require.Equal(t, pickupResponse2.Book.Authors[0].AuthorId, book2.Book.Authors[0].AuthorId)
	require.Equal(t, pickupResponse2.Book.Authors[0].Name, book2.Book.Authors[0].Name)

}