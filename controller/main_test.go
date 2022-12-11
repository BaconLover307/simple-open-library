package controller_test

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"simple-open-library/test"
	"testing"
)

var (
	testDB *sql.DB
	testCtx context.Context
	testRouter http.Handler
	BaseURL string = "http://" + os.Getenv("SV_TEST_HOST") + ":" + os.Getenv("SV_TEST_PORT")
)

func TestMain(m *testing.M) {
	testDB = test.SetupTestDB()
	test.TruncateDatabase(testDB)
	testRouter = test.InitializeTestServer(testDB)

	os.Exit(m.Run())
}