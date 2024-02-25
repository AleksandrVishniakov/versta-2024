package handlers

import (
	"log/slog"
	"net/http"
	"time"
)

type writerRecorder struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (r *writerRecorder) WriteHeader(status int) {
	r.status = status
}

func (r *writerRecorder) Write(bytes []byte) (int, error) {
	r.body = bytes

	return r.ResponseWriter.Write(bytes)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var startedAt = time.Now()

		recorder := &writerRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
			body:           make([]byte, 0),
		}

		next.ServeHTTP(recorder, r)

		var duration = time.Since(startedAt)
		var url = r.URL.String()
		var statusCode = recorder.status
		var method = r.Method

		log := loggingMethod(statusCode)

		log("http",
			slog.Group("request",
				slog.String("url", url),
				slog.String("method", method),
			),

			slog.Group("response",
				slog.Int("code", statusCode),
				slog.String("duration", duration.String()),
			),
		)
	})
}

func loggingMethod(status int) func(msg string, args ...any) {
	if status >= 200 && status < 300 {
		return slog.Info
	} else {
		return slog.Warn
	}
}
