package exception

import (
	"net/http"
	"simple-open-library/helper"
	"simple-open-library/model/web"

	"github.com/go-playground/validator/v10"
)

func HttprouterErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
	}
	if conflictError(writer, request, err) {
		return
	}
	if badRequestError(writer, request, err) {
		return
	}
	if validationErrors(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func ChiErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			err := recover()

			if notFoundError(writer, request, err) {
				return
			}
			if conflictError(writer, request, err) {
				return
			}
			if badRequestError(writer, request, err) {
				return
			}
			if validationErrors(writer, request, err) {
				return
			}
			if internalServerError(writer, request, err) {
				return
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Data:   exception.Error(),
		}

		helper.WriteResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func conflictError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(ConflictError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)

		webResponse := web.WebResponse{
			Code:   http.StatusConflict,
			Status: "CONFLICT",
			Data:   exception.Error(),
		}

		helper.WriteResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(BadRequestError)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}

		helper.WriteResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data:   exception.Error(),
		}

		helper.WriteResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)

		webResponse := web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Data:   exception.Error(),
		}

		helper.WriteResponseBody(writer, webResponse)
		return true
	} else {
		return false
	}
}
