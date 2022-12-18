package controller

import (
	"net/http"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PickupController interface {
	ListSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetScheduleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	SubmitSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}

type pickupController struct {
	PickupService service.PickupService
	BookService   service.BookService
}

func NewPickupController(pickupService service.PickupService, bookService service.BookService) PickupController {
	return &pickupController{
		PickupService: pickupService,
		BookService:   bookService,
	}
}

func (controller pickupController) ListSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupResponses := controller.PickupService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pickupResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller pickupController) GetScheduleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupId := params.ByName("pickupId")
	id, err := strconv.Atoi(pickupId)
	helper.PanicIfError(err)

	pickupResponse := controller.PickupService.FindById(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pickupResponse,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller pickupController) SubmitSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupCreateRequest := web.PickupCreateRequest{}
	err := helper.ReadFromRequestBody(request, &pickupCreateRequest)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	controller.BookService.SaveBook(request.Context(), pickupCreateRequest.Book)

	submitResponse := controller.PickupService.Create(request.Context(), pickupCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   submitResponse,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller pickupController) UpdateSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupUpdateScheduleRequest := web.PickupUpdateScheduleRequest{}
	err := helper.ReadFromRequestBody(request, &pickupUpdateScheduleRequest)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	pickupId := params.ByName("pickupId")
	id, err := strconv.Atoi(pickupId)
	helper.PanicIfError(err)
	pickupUpdateScheduleRequest.PickupId = id

	updateResponse := controller.PickupService.UpdateSchedule(request.Context(), pickupUpdateScheduleRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   updateResponse,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller pickupController) DeleteSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupId := params.ByName("pickupId")
	id, err := strconv.Atoi(pickupId)
	helper.PanicIfError(err)

	controller.PickupService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteResponseBody(writer, webResponse)
}
