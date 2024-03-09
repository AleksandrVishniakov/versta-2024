package handlers

import (
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/sessionsservice"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/usersservice"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/parser"
)

const (
	codeQueryName           = "code"
	emailQueryName          = "email"
	sessionKeyQueryName     = "session_key"
	needToSendMailQueryName = "send_email"

	codeLength = 6
)

type HTTPHandler struct {
	userService     usersservice.UsersService
	sessionsService sessionsservice.SessionsService
	cookieTTL       time.Duration
}

func NewHTTPHandler(
	userService usersservice.UsersService,
	sessionsService sessionsservice.SessionsService,
	cookieTTL time.Duration,
) *HTTPHandler {
	return &HTTPHandler{
		userService:     userService,
		sessionsService: sessionsService,
		cookieTTL:       cookieTTL,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.HandleFunc("GET /api/auth", h.handleAuthentication)
	mux.HandleFunc("GET /api/{email}/verify", h.handleEmailVerification)

	mux.HandleFunc("GET /api/user", h.getUser)
	mux.HandleFunc("PUT /api/{email}/name", h.updateName)

	mux.HandleFunc("GET /api/internal/user/{email}/verification_code", h.getUserVerificationCode)

	return Recovery(Logger(CORS(Cookies(mux))))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) handleAuthentication(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get(emailQueryName)
	var needToSendEmail = true
	var err error

	if needToSendEmailStr := r.URL.Query().Get(needToSendMailQueryName); needToSendEmailStr != "" {
		needToSendEmail, err = strconv.ParseBool(needToSendEmailStr)
		if err != nil {
			e.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := h.userService.Register(email, needToSendEmail)
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
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if errors.Is(err, usersservice.ErrMismatchedCodes) {
		e.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !isCookiesAccepted(r) {
		return
	}

	user, err := h.userService.FindByEmail(email)
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionKey, err := h.sessionsService.Create(user.Id)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := sessionCookie(sessionKey, h.cookieTTL)
	http.SetCookie(w, cookie)
}

func (h *HTTPHandler) getUser(w http.ResponseWriter, r *http.Request) {
	var email = r.URL.Query().Get(emailQueryName)
	var sessionKey = r.URL.Query().Get(sessionKeyQueryName)

	var user = &usersservice.UserResponseDTO{}
	var err error

	if email != "" {
		user, err = h.userService.FindByEmail(email)
	} else if sessionKey != "" {
		user, sessionKey, err = h.findBySessionKey(sessionKey)
		if errors.Is(err, sessionsservice.ErrSessionNotFound) {
			e.WriteError(w, http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, sessionsservice.ErrSessionExpired) {
			e.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err != nil {
			e.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if isCookiesAccepted(r) {
			http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))
		}
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
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

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

func (h *HTTPHandler) findBySessionKey(sessionKey string) (*usersservice.UserResponseDTO, string, error) {
	err := h.sessionsService.Valid(sessionKey)
	if err != nil {
		return nil, "", err
	}

	sessionKey, err = h.sessionsService.UpdateKey(sessionKey)
	if err != nil {
		return nil, "", err
	}

	user, err := h.userService.FindBySessionKey(sessionKey)
	if err != nil {
		return nil, "", err
	}

	return user, sessionKey, nil
}

func (h *HTTPHandler) updateName(w http.ResponseWriter, r *http.Request) {
	type nameDTO struct {
		Name string `json:"name"`
	}

	email := r.PathValue(emailQueryName)

	userName, err := parser.Decode[nameDTO](r.Body)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(userName.Name) < 2 {
		e.WriteError(w, http.StatusBadRequest, "name should contain at least 2 symbols")
		return
	}

	user, err := h.userService.FindByEmail(email)
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.userService.UpdateName(user.Id, userName.Name)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func isCookiesAccepted(r *http.Request) bool {
	isAccepted, ok := (r.Context().Value(IsCookieAcceptedKey)).(bool)
	if !ok {
		return false
	}

	return isAccepted
}

func sessionCookie(sessionKey string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     "sessionKey",
		Value:    sessionKey,
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}
}
