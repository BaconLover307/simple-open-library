//go:build wireinject
// +build wireinject


package app

import (
	"net/http"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/repository"
	"simple-open-library/service"
	"simple-open-library/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
)

var librarySet = wire.NewSet(
	lib.NewOpenLibraryLib,
)

var repositorySet = wire.NewSet(
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
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		middleware.NewAuthMiddleware,
		NewServer,
	)
	return nil
}