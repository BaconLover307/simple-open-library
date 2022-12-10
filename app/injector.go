//go:build wireinject
// +build wireinject

package app

import (
	"net/http"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/repository"
	"simple-open-library/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

var librarySet = wire.NewSet(
	lib.NewOpenLibraryLib,
)

var repositorySet = wire.NewSet(
	repository.NewOpenLibraryRepository,
	repository.NewBookRepository,
	repository.NewPickupRepository,
)

var serviceSet = wire.NewSet(
	service.NewBookService,
	service.NewLibraryService,
	service.NewPickupService,
)

var controllerSet = wire.NewSet(
	controller.NewLibraryController,
	controller.NewPickupController,
)

func InitializeServer() *http.Server {
	wire.Build(
		NewDB,
		validator.New,
		librarySet,
		repositorySet,
		serviceSet,
		controllerSet,
		NewRouter,
		NewServer,
	)
	return nil
}