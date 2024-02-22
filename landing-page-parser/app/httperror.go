package main

import (
	"log/slog"
	"net/http"
	"time"
)

type HandlerWithErr func(w http.ResponseWriter, r *http.Request) error

type ResponseError struct {
	Code          int       `json:"code"`
	Message       string    `json:"message"`
	Timestamp     time.Time `json:"timestamp"`
	DeveloperCode int       `json:"developerCode"`
}

func ErrorHandler(next HandlerWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)

		if err == nil {
			return
		}

		var msg = err.Error()

		respError := ResponseError{
			Code:          http.StatusInternalServerError,
			Message:       msg,
			Timestamp:     time.Now(),
			DeveloperCode: http.StatusInternalServerError * 1000,
		}

		w.Header().Set("Content-Type", "application/json")
		err = Encode(w, respError)
		if err != nil {
			slog.Error("response error encoding error",
				slog.String("error", err.Error()),
			)
		}
		w.WriteHeader(http.StatusInternalServerError)

		slog.Debug("http_error",
			slog.Int("code", http.StatusInternalServerError),
			slog.String("message", msg),
			slog.String("time", time.Now().String()),
			slog.Int("developer_code", 500000),
		)
	}
}
