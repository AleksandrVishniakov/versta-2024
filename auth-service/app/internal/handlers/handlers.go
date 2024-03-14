package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/services/usersservice"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/jwttokens"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/parser"
)

const (
	codeQueryName           = "code"
	emailQueryName          = "email"
	needToSendMailQueryName = "send_email"

	codeLength = 6
)

type HTTPHandler struct {
	userService     usersservice.UsersService
	tokensManager   jwttokens.Manager
	refreshTokenTTL time.Duration
}

func NewHTTPHandler(
	userService usersservice.UsersService,
	tokensManager jwttokens.Manager,
	cookieTTL time.Duration,
) *HTTPHandler {
	return &HTTPHandler{
		userService:     userService,
		tokensManager:   tokensManager,
		refreshTokenTTL: cookieTTL,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	jwtAuth := NewJWTAuthMiddleware(h.tokensManager)

	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.HandleFunc("GET /api/auth", h.handleAuthentication)
	mux.HandleFunc("GET /api/tokens/refresh", h.refreshTokens)
	mux.HandleFunc("GET /api/{email}/verify", h.handleEmailVerification)

	mux.HandleFunc("GET /api/user/email/{email}", h.getUserByEmail)

	mux.Handle("PUT /api/user/name", jwtAuth(http.HandlerFunc(h.updateName)))

	mux.Handle("GET /api/user/my_profile", jwtAuth(http.HandlerFunc(h.getUserByToken)))

	mux.HandleFunc("GET /api/internal/user/{email}/verification_code", h.getUserVerificationCode)

	return Recovery(Logger(CORS(mux)))
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

	user, err := h.userService.FindByEmail(email)
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.tokensManager.CreateTokens(user.Id, jwttokens.AccessTokenPayload{
		UserId: user.Id,
		Email:  user.Email,
	})
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, getRefreshTokenCookie(tokens.RefreshToken, h.refreshTokenTTL))

	err = parser.EncodeResponse(w, tokens.AccessToken, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	var email = r.PathValue(emailQueryName)

	var user = &usersservice.UserResponseDTO{}
	var err error

	if email == "" {
		e.WriteError(w, http.StatusBadRequest, "no email provided")
		return
	}

	user, err = h.userService.FindByEmail(email)
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = parser.EncodeResponse(w, user, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) getUserByToken(w http.ResponseWriter, r *http.Request) {
	email := fmt.Sprintf("%v", r.Context().Value(EmailContextKey))

	user, err := h.userService.FindByEmail(email)
	if errors.Is(err, usersservice.ErrUserNotFound) {
		e.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
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

func (h *HTTPHandler) updateName(w http.ResponseWriter, r *http.Request) {
	type nameDTO struct {
		Name string `json:"name"`
	}

	email := fmt.Sprintf("%v", r.Context().Value(EmailContextKey))

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

func (h *HTTPHandler) refreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie(RefreshTokenCookieKey)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	refreshTokenPayload, err := h.tokensManager.ParseRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		e.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.userService.FindById(refreshTokenPayload.UserId)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokens, err := h.tokensManager.RefreshTokens(refreshTokenCookie.Value, &jwttokens.AccessTokenPayload{
		UserId: user.Id,
		Email:  user.Email,
	})

	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, getRefreshTokenCookie(tokens.RefreshToken, h.refreshTokenTTL))

	err = parser.EncodeResponse(w, tokens.AccessToken, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func getRefreshTokenCookie(token string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieKey,
		Value:    token,
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
		HttpOnly: true,
	}
}
