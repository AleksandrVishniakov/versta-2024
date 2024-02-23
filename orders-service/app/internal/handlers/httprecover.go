package handlers

import (
	"log/slog"
	"net/http"
)

func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				slog.Error("panic recovered",
					slog.String("error", err.(string)),

					slog.Group("request",
						slog.String("url", r.URL.String()),
						slog.String("method", r.Method),
					),
				)
			}
		}()

		next.ServeHTTP(w, r)
	}
}
