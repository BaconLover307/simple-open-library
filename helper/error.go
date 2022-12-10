package helper

import "log"

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