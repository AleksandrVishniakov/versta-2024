package e

import (
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/parser"
	"log/slog"
	"net/http"
	"time"
)

type ResponseError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func WriteError(w http.ResponseWriter, code int, message string) {
	rError := ResponseError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}

	err := parser.EncodeResponse(w, rError, code)
	if err != nil {
		slog.Error(
			"error parsing response error",
			slog.String("error", err.Error()),
		)
	}
}
