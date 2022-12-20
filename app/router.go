package app

import (
	"simple-open-library/controller"
	"simple-open-library/exception"
	mid "simple-open-library/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(pickupController controller.PickupController, libraryController controller.LibraryController, bookController controller.BookController) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(exception.ChiErrorHandler)

	// $ Need API Key
	router.Group(func(router chi.Router) {
		router.Use(mid.ChiAuthMiddleware)
		router.Get("/api/pickups", pickupController.ListSchedule)
		router.Get("/api/pickups/{pickupId}", pickupController.GetScheduleById)
		router.Post("/api/pickups", pickupController.SubmitSchedule)
		router.Put("/api/pickups/{pickupId}", pickupController.UpdateSchedule)
		router.Delete("/api/pickups/{pickupId}", pickupController.DeleteSchedule)
	})

	// $ Public API
	router.Group(func(router chi.Router) {
		router.Get("/api/subjects/{subject}", libraryController.BrowseBySubject)
		router.Get("/api/books", bookController.ListBooks)
	})

	return router
}

// func NewRouteExclusions() *route.Prefixes {
// 	prefixes := route.NewPrefixes()
// 	prefixes.Add("/api/books")
// 	prefixes.Add("/api/subjects")

// 	return prefixes
// }
