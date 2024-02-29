package handlers

import (
	"io"
	"log/slog"
	"net/http"
)

const (
	codeQueryName       = "code"
	emailQueryName      = "email"
	sessionKeyQueryName = "sessionKey"
)

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.HandleFunc("GET /api/auth", h.handleAuthentication)
	mux.HandleFunc("GET /api/{email}/verify", h.handleEmailVerification)

	mux.HandleFunc("GET /api/user", h.getUser)

	mux.HandleFunc("GET /api/internal/user/{email}/verification_code", h.getUserVerificationCode)

	return Recovery(Logger(Cookies(mux)))
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

func (h *HTTPHandler) getUserVerificationCode(w http.ResponseWriter, r *http.Request) {

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
