package controller

import (
	"net/http"
	"simple-open-library/exception"
	"simple-open-library/helper"
	"simple-open-library/model/web"
	"simple-open-library/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PickupController interface {
	ListSchedule(writer http.ResponseWriter, request *http.Request)
	GetScheduleById(writer http.ResponseWriter, request *http.Request)
	SubmitSchedule(writer http.ResponseWriter, request *http.Request)
	UpdateSchedule(writer http.ResponseWriter, request *http.Request)
	DeleteSchedule(writer http.ResponseWriter, request *http.Request)
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

func (controller pickupController) ListSchedule(writer http.ResponseWriter, request *http.Request) {
	pickupResponses := controller.PickupService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   pickupResponses,
	}

	helper.WriteResponseBody(writer, webResponse)
}

func (controller pickupController) GetScheduleById(writer http.ResponseWriter, request *http.Request) {
	pickupId := chi.URLParam(request, "pickupId")
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

func (controller pickupController) SubmitSchedule(writer http.ResponseWriter, request *http.Request) {
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

func (controller pickupController) UpdateSchedule(writer http.ResponseWriter, request *http.Request) {
	pickupUpdateScheduleRequest := web.PickupUpdateScheduleRequest{}
	err := helper.ReadFromRequestBody(request, &pickupUpdateScheduleRequest)
	if err != nil {
		panic(exception.NewBadRequestError(err.Error()))
	}

	pickupId := chi.URLParam(request, "pickupId")
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

func (controller pickupController) DeleteSchedule(writer http.ResponseWriter, request *http.Request) {
	pickupId := chi.URLParam(request, "pickupId")
	id, err := strconv.Atoi(pickupId)
	helper.PanicIfError(err)

	controller.PickupService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
	}

	helper.WriteResponseBody(writer, webResponse)
}
