// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"net/http"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/middleware"
	"simple-open-library/repository"
	"simple-open-library/service"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from injector.go:

func InitializeServer() *http.Server {
	pickupRepository := repository.NewPickupRepository()
	db := NewDB()
	validate := validator.New()
	pickupService := service.NewPickupService(pickupRepository, db, validate)
	bookRepository := repository.NewBookRepository()
	bookService := service.NewBookService(bookRepository, db, validate)
	pickupController := controller.NewPickupController(pickupService, bookService)
	openLibraryLib := lib.NewOpenLibraryLib()
	libraryService := service.NewLibraryService(openLibraryLib, db, validate)
	libraryController := controller.NewLibraryController(libraryService)
	bookController := controller.NewBookController(bookService)
	router := NewRouter(pickupController, libraryController, bookController)
	authMiddleware := middleware.NewAuthMiddleware(router)
	server := NewServer(authMiddleware)
	return server
}

// injector.go:

var librarySet = wire.NewSet(lib.NewOpenLibraryLib)

var repositorySet = wire.NewSet(repository.NewBookRepository, repository.NewPickupRepository)

var serviceSet = wire.NewSet(service.NewBookService, service.NewLibraryService, service.NewPickupService)

var controllerSet = wire.NewSet(controller.NewLibraryController, controller.NewPickupController, controller.NewBookController)
