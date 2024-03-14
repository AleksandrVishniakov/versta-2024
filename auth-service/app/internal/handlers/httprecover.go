package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				slog.Error("panic recovered",
					slog.String("error", fmt.Sprintf("%v", err)),

					slog.Group("request",
						slog.String("url", r.URL.String()),
						slog.String("method", r.Method),
					),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
