package main

import (
	"simple-open-library/app"
	"simple-open-library/helper"
)

func main() {
	server := app.InitializeServer()
	err := server.ListenAndServe()
	helper.PanicIfError(err)
}