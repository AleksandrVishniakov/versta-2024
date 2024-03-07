package handlers

import (
	"errors"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/e"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request) (statusCode int, err error)

func Errors(next ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusCode, err := next(w, r)
		if err != nil {
			var responseError *e.ResponseError
			switch {
			case errors.As(err, &responseError):
				e.WriteError(w, responseError.Code, responseError.Message)
				return
			default:
				e.WriteError(w, statusCode, err.Error())
				return
			}
		}
	})
}
