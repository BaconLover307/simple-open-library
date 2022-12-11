package repository_test

import (
	"database/sql"
	"os"
	"simple-open-library/test"
	"testing"
)

var (
	testDB *sql.DB
	testTx *sql.Tx
)

func TestMain(m *testing.M) {
	testDB = test.SetupTestDB()
	test.TruncateDatabase(testDB)	

	os.Exit(m.Run())
}


