package handlers

import (
	"io"
	"log/slog"
	"net/http"
)

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.HandleFunc("GET /auth", h.handleAuthentication)
	mux.HandleFunc("GET /{email}/verify", h.handleEmailVerification)

	mux.HandleFunc("GET /user/{id}", h.getUser)

	return Recovery(Logger(mux))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleAuthentication(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) handleEmailVerification(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) getUser(w http.ResponseWriter, r *http.Request) {

}

func closeReadCloser(r io.ReadCloser) {
	err := r.Close()
	if err != nil {
		slog.Error(
			"read closer closing error",
			slog.String("error", err.Error()),
		)
	}
}
