// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package test

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"simple-open-library/app"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/repository"
	"simple-open-library/service"
)

// Injectors from test_injector.go:

func InitializeTestRouter(db *sql.DB) *chi.Mux {
	pickupRepository := repository.NewPickupRepository()
	validate := validator.New()
	pickupService := service.NewPickupService(pickupRepository, db, validate)
	bookRepository := repository.NewBookRepository()
	bookService := service.NewBookService(bookRepository, db, validate)
	pickupController := controller.NewPickupController(pickupService, bookService)
	openLibraryLib := lib.NewOpenLibraryLib()
	libraryService := service.NewOpenLibraryService(openLibraryLib, validate)
	libraryController := controller.NewLibraryController(libraryService)
	bookController := controller.NewBookController(bookService)
	mux := app.NewRouter(pickupController, libraryController, bookController)
	return mux
}

// test_injector.go:

var librarySet = wire.NewSet(lib.NewOpenLibraryLib)

var repositorySet = wire.NewSet(repository.NewBookRepository, repository.NewPickupRepository)

var serviceSet = wire.NewSet(service.NewBookService, service.NewOpenLibraryService, service.NewPickupService)

var controllerSet = wire.NewSet(controller.NewLibraryController, controller.NewPickupController, controller.NewBookController)
