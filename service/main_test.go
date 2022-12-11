package service_test

import (
	"context"
	"database/sql"
	"os"
	"simple-open-library/lib"
	"simple-open-library/repository"
	"simple-open-library/service"
	"simple-open-library/test"
	"testing"

	"github.com/go-playground/validator/v10"
)

var (
	testDB *sql.DB
	testLibraryService service.LibraryService
	testBookService service.BookService
	testPickupService service.PickupService
	testCtx context.Context
)

func TestMain(m *testing.M) {
	testDB = test.SetupTestDB()
	test.TruncateDatabase(testDB)
	validate := validator.New()
	testCtx = context.Background()

	testLibraryService = service.NewLibraryService(lib.NewOpenLibraryLib(), testDB, validate)
	testBookService = service.NewBookService(repository.NewBookRepository(), testDB, validate)
	testPickupService = service.NewPickupService(repository.NewPickupRepository(), testDB, validate)

	os.Exit(m.Run())
}