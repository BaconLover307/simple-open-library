package controller

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type PickupControllerImpl struct {
	PickupService service.PickupService
	BookService service.BookService 
}

func NewPickupController(pickupService service.PickupService, bookService service.BookService) PickupController {
	return &PickupControllerImpl{
		PickupService: pickupService,
		BookService: bookService,
	}
}

func (controller PickupControllerImpl) ListSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupResponses := controller.PickupService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pickupResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller PickupControllerImpl) GetScheduleById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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

func (controller PickupControllerImpl) SubmitSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupCreateRequest := web.PickupCreateRequest{}
	helper.ReadFromRequestBody(request, &pickupCreateRequest)

	controller.BookService.SaveBook(request.Context(), pickupCreateRequest.Book)

	submitResponse := controller.PickupService.Create(request.Context(), pickupCreateRequest)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   submitResponse,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller PickupControllerImpl) UpdateSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	pickupUpdateScheduleRequest := web.PickupUpdateScheduleRequest{}
	helper.ReadFromRequestBody(request, &pickupUpdateScheduleRequest)
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

func (controller PickupControllerImpl) DeleteSchedule(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
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
