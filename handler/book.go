package handler

import (
	error2 "github.com/ShaDream/library-api-practice/model/error"
	"github.com/ShaDream/library-api-practice/service"
	"github.com/gorilla/mux"
	"strconv"
)

func initBookRoutes(r *mux.Router) {

	HandleFunc(r, "/api/book/{id}", GetBook).
		Methods("GET")

	HandleFunc(r, "/api/book/{id}/items", GetPhysicalBooks).
		Methods("GET")

	HandleFunc(r, "/api/book", GetBooks).
		Methods("GET").
		Queries("page", "{page:[0-9]+}", "size", "{size:[0-9]+}")
}

func GetBooks(service service.Service, params map[string]string) Result {
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		return GetResultError(error2.NewUnprocessableEntityError(err))
	}
	size, err := strconv.Atoi(params["size"])
	if err != nil {
		return GetResultError(error2.NewUnprocessableEntityError(err))
	}

	books, err := service.GetBooks(page, size)
	if err != nil {
		return GetResultError(err)
	}
	return Result{
		Data: books,
	}
}

func GetBook(service service.Service, params map[string]string) Result {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return GetResultError(error2.NewUnprocessableEntityError(err))
	}
	book, err := service.GetBook(id)
	if err != nil {
		return GetResultError(err)
	}
	return Result{
		Data: book,
	}
}

func GetPhysicalBooks(service service.Service, params map[string]string) Result {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return GetResultError(error2.NewUnprocessableEntityError(err))
	}
	books, err := service.GetPhysicalBooks(id)
	if err != nil {
		return GetResultError(err)
	}
	return Result{
		Data: books,
	}
}

func GetResultError(err error) Result {
	return Result{
		Error: err,
	}
}
