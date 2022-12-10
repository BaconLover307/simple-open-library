package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PickupController interface {
	ListSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetScheduleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SubmitSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}