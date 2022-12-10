package app

import (
	"simple-open-library/controller"
	"simple-open-library/exception"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(pickupController controller.PickupController, libraryController controller.LibraryController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/pickups", pickupController.ListSchedule)
	router.GET("/api/pickups/:pickupId", pickupController.GetScheduleById)
	router.POST("/api/pickups", pickupController.SubmitSchedule)
	router.PUT("/api/pickups/:pickupId", pickupController.UpdateSchedule)
	router.DELETE("/api/pickups/:pickupId", pickupController.DeleteSchedule)
	
	router.GET("/api/subjects/:subject", libraryController.BrowseBySubject)
	
	router.PanicHandler = exception.ErrorHandler

	return router
}