package repository_test

import (
	"context"
	"simple-open-library/helper"
	"simple-open-library/model/domain"
	"simple-open-library/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	inputPickup1 = domain.Pickup{
		Book: inputBook1,
		Schedule: "2021-01-01 10:10:10",
	}
	inputPickup2 = domain.Pickup{
		Book: inputBook2,
		Schedule: "2022-02-02 12:20:20",
	}
	inputPickup3 = domain.Pickup{
		Book: inputBook2,
		Schedule: "2013-01-03 13:23:33",
	}
	inputSchedule1 = "2003-03-03 03:30:30"
)

func TestRepoPickupCreate(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	savedPickup := pickupRepo.Create(ctx, testTx, inputPickup1)
	pickupResult, err := pickupRepo.FindById(ctx, testTx, savedPickup.PickupId)
	inputPickup1.PickupId = savedPickup.PickupId
	
	require.NoError(t, err)
	require.Equal(t, inputPickup1, pickupResult)
}

func TestRepoPickupFindByIdFail(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupResult, err := pickupRepo.FindById(ctx, testTx, inputPickup2.PickupId)
	
	require.Error(t, err)
	require.Empty(t, pickupResult)
}

func TestRepoPickupUpdate(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)
	
	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	savedPickup := pickupRepo.Create(ctx, testTx, inputPickup2)
	pickupResult := pickupRepo.UpdateSchedule(ctx, testTx, savedPickup)
	inputPickup2.PickupId = savedPickup.PickupId
	
	require.NoError(t, err)
	require.Equal(t, inputPickup2, pickupResult)
}


func TestRepoPickupDeleteSuccess(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Delete(ctx, testTx, inputPickup2)
	pickupResult, err := pickupRepo.FindById(ctx, testTx, inputPickup2.PickupId)
	
	require.Error(t, err)
	require.Empty(t, pickupResult)
}

func TestRepoPickupDeleteFail(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	pickupRepo.Delete(ctx, testTx, inputPickup2)
	pickupResult, err := pickupRepo.FindById(ctx, testTx, inputPickup2.PickupId)
	
	require.Error(t, err)
	require.Empty(t, pickupResult)
}

func TestRepoPickupFindAll(t *testing.T) {
	testTx, err := testDB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(testTx)

	ctx := context.Background()
	pickupRepo := repository.NewPickupRepository()
	savedPickup := pickupRepo.Create(ctx, testTx, inputPickup3)
	inputPickup3.PickupId = savedPickup.PickupId

	pickupsResult := pickupRepo.FindAll(ctx, testTx)

	require.Len(t, pickupsResult, 2)
	require.Equal(t, inputPickup1, pickupsResult[0])
	require.Equal(t, inputPickup3, pickupsResult[1])
}