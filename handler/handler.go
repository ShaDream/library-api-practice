package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ShaDream/library-api-practice/database"
	error2 "github.com/ShaDream/library-api-practice/model/error"
	"github.com/ShaDream/library-api-practice/repository"
	"github.com/ShaDream/library-api-practice/service"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GetRouter() *mux.Router {
	apiRoute := mux.NewRouter()

	apiRoute.Use(LoggingMiddleware)
	apiRoute.Use(GetTimeoutMiddleware(time.Second * 60))
	initBookRoutes(apiRoute)
	return apiRoute
}

type Func func(service service.Service, params map[string]string) Result

type Result struct {
	Data  interface{}
	Error error
}

func HandleFunc(r *mux.Router, path string, f Func) *mux.Route {
	return r.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		tx := database.GetTransactionWithContext(ctx)
		defer func() {
			if p := recover(); p != nil {
				tx.Rollback()
			}
		}()
		s := service.NewService(repository.NewRepository(tx))
		params := mux.Vars(request)
		resultChan := make(chan Result)

		go func(r chan<- Result) {
			r <- f(*s, params)
		}(resultChan)

		select {
		case <-ctx.Done():
			fmt.Println("test")
			writer.WriteHeader(http.StatusGatewayTimeout)
			tx.Rollback()
			return
		case res := <-resultChan:
			hasError := handleResult(res, writer)
			if hasError {
				tx.Rollback()
			} else {
				tx.Commit()
			}
			return
		}
	})
}

// returns true if in result has error or error occurs when writing result
func handleResult(result Result, writer http.ResponseWriter) bool {
	if result.Error != nil {
		var statusError error2.StatusError
		switch t := result.Error.(type) {
		case error2.StatusError:
			statusError = t
		default:
			statusError = error2.NewInternalServerError(t)
		}
		writer.WriteHeader(statusError.StatusCode)
		writer.Write([]byte(result.Error.Error()))
		return true
	} else {
		dataMarshal, err := json.Marshal(result.Data)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return true
		}
		_, err = writer.Write(dataMarshal)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return true
		}
		return false
	}
}
