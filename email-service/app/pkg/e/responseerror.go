package e

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/email-service/app/pkg/parser"
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
