package controller_test

import (
	"os"
	"testing"
)

var (
	BaseURL string = "http://" + os.Getenv("SV_TEST_HOST") + ":" + os.Getenv("SV_TEST_PORT")
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}