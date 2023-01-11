//go:build wireinject
// +build wireinject

package test

import (
	"database/sql"
	"simple-open-library/app"
	"simple-open-library/controller"
	"simple-open-library/lib"
	"simple-open-library/repository"
	"simple-open-library/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
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
	service.NewOpenLibraryService,
	service.NewPickupService,
)

var controllerSet = wire.NewSet(
	controller.NewLibraryController,
	controller.NewPickupController,
	controller.NewBookController,
)

func InitializeTestRouter(db *sql.DB) *chi.Mux {
	wire.Build(
		validator.New,
		librarySet,
		repositorySet,
		serviceSet,
		controllerSet,
		app.NewRouter,
	)
	return nil
}
