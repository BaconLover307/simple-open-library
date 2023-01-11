package helper

import (
	"log"
	"testing"
)

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func FatalIfError(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}

func FatalIfMockError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
}
