package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/usersservice"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/parser"
	"io"
	"log/slog"
	"net/http"
)

const (
	codeQueryName       = "code"
	emailQueryName      = "email"
	sessionKeyQueryName = "sessionKey"

	codeLength = 6
)

type HTTPHandler struct {
	userService usersservice.UsersService
}

func NewHTTPHandler(userService usersservice.UsersService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
	}
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
	email := r.URL.Query().Get(emailQueryName)

	id, err := h.userService.Register(email)
	if err != nil {
		e.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = parser.EncodeResponse(w, id, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) handleEmailVerification(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue(emailQueryName)
	code := r.URL.Query().Get(codeQueryName)

	if len(code) != codeLength {
		e.WriteError(w, http.StatusBadRequest, "invalid code length")
		return
	}

	err := h.userService.VerifyEmail(email, code)
	if errors.Is(err, usersservice.ErrMismatchedCodes) {
		e.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) getUser(w http.ResponseWriter, r *http.Request) {
	var email = r.URL.Query().Get(emailQueryName)
	var sessionKey = r.URL.Query().Get(sessionKeyQueryName)

	var user = &usersservice.UserResponseDTO{}
	var err error

	if email != "" {
		user, err = h.userService.FindByEmail(email)
	} else if sessionKey != "" {
		user, err = h.userService.FindBySessionKey(sessionKey)
	} else {
		e.WriteError(w, http.StatusBadRequest, "no email or sessionKey parameter provided")
		return
	}

	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = parser.EncodeResponse(w, user, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) getUserVerificationCode(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue(emailQueryName)

	verificationCode, err := h.userService.GetVerificationCode(email)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = parser.EncodeResponse(w, verificationCode, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
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
