package handlers

import (
	"errors"
	"net/http"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
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
