package app

import (
	"simple-open-library/controller"
	"simple-open-library/exception"
	"simple-open-library/model/route"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(pickupController controller.PickupController, libraryController controller.LibraryController, bookController controller.BookController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/pickups", pickupController.ListSchedule)
	router.GET("/api/pickups/:pickupId", pickupController.GetScheduleById)
	router.POST("/api/pickups", pickupController.SubmitSchedule)
	router.PUT("/api/pickups/:pickupId", pickupController.UpdateSchedule)
	router.DELETE("/api/pickups/:pickupId", pickupController.DeleteSchedule)
	
	router.GET("/api/subjects/:subject", libraryController.BrowseBySubject)
	
	router.GET("/api/books", bookController.ListBooks)
	
	router.PanicHandler = exception.ErrorHandler

	return router
}

func NewRouteExclusions() *route.Prefixes {
	prefixes := route.NewPrefixes()
	prefixes.Add("/api/books")
	prefixes.Add("/api/subjects")

	return prefixes
}