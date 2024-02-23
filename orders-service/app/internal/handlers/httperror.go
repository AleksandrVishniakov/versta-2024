package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/parser"
)

type HandlerWithErr func(w http.ResponseWriter, r *http.Request) error

type ResponseError struct {
	Code          int       `json:"code"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
	DeveloperCode int       `json:"developerCode"`
}

func Errors(next HandlerWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)

		if err == nil {
			return
		}

		var msg = err.Error()

		devCode, exists := errorCodes[err]
		if !exists {
			devCode = http.StatusInternalServerError * 500000
		}

		statusCode := devCode / 1000

		respError := ResponseError{
			Code:          statusCode,
			Message:       msg,
			Timestamp:     time.Now(),
			DeveloperCode: devCode,
		}

		w.Header().Set("Content-Type", "application/json")
		err = parser.Encode(w, respError)
		if err != nil {
			slog.Error("response error encoding error",
				slog.String("error", err.Error()),
			)
		}
		w.WriteHeader(http.StatusInternalServerError)

		slog.Debug("http_error",
			slog.Int("code", statusCode),
			slog.String("message", msg),
			slog.String("time", time.Now().String()),
			slog.Int("developer_code", devCode),
		)
	}
}

var errorCodes = map[error]int{}
