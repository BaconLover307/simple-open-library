//go:build wireinject
// +build wireinject

package test

import (
	"database/sql"
	"net/http"
	"simple-open-library/app"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/middleware"
	"simple-open-library/repository"
	"simple-open-library/service"

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
	controller.NewBookController,
)

func InitializeTestServer(db *sql.DB) *middleware.AuthMiddleware {
	wire.Build(
		validator.New,
		librarySet,
		repositorySet,
		serviceSet,
		controllerSet,
		app.NewRouter,
		wire.Bind(new(http.Handler), new(*httprouter.Router)),
		app.NewRouteExclusions,
		middleware.NewAuthMiddleware,
	)
	return nil
}