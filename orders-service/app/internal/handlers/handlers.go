package handlers

import (
	"errors"
	"fmt"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/api/authapi"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/api/emailapi"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/orders"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/ordersservice"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/parser"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/utils"
)

const (
	sessionKeyLength = 32

	emailQueryName            = "email"
	verificationCodeQueryName = "code"
)

type HTTPHandler struct {
	ordersService ordersservice.OrdersService
	authAPI       authapi.API
	emailAPI      emailapi.API
	cookieTTL     time.Duration
}

func NewHTTPHandler(
	ordersService ordersservice.OrdersService,
	authAPI authapi.API,
	emailAPI emailapi.API,
	cookieTTL time.Duration,
) *HTTPHandler {
	return &HTTPHandler{
		ordersService: ordersService,
		authAPI:       authAPI,
		emailAPI:      emailAPI,
		cookieTTL:     cookieTTL,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.Handle("GET /api/orders", Errors(h.getAllOrders))
	mux.Handle("POST /api/order", Errors(h.createNewOrder))
	mux.Handle("GET /api/order/{orderId}", Errors(h.getOrder))
	mux.Handle("DELETE /api/order/{orderId}", Errors(h.deleteOrder))
	mux.Handle("GET /api/order/{orderId}/verify", Errors(h.verifyOrder))
	mux.Handle("GET /api/order/{orderId}/complete", Errors(h.completeOrder))

	return Recovery(Logger(Cookies(mux)))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) getAllOrders(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))
	if len(sessionKey) != sessionKeyLength {
		return http.StatusUnauthorized, errors.New("invalid session key")
	}

	user, sessionKey, err := h.authAPI.FindBySessionKey(sessionKey)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	userOrders, err := h.ordersService.FindAll(user.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, userOrders, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	return 0, nil
}

func (h *HTTPHandler) createNewOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)
	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))

	email := r.URL.Query().Get(emailQueryName)

	var userId int
	var err error
	var isNewUser = false

	if len(sessionKey) != sessionKeyLength {
		userId, err = h.authAPI.Create(email, false)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		isNewUser = true
	} else {
		user, sessionKey, err := h.authAPI.FindBySessionKey(sessionKey)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		userId = user.Id

		email = user.Email

		http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))
	}

	var emailVerificationCode string
	if isNewUser {
		verificationCode, err := h.authAPI.GetVerificationCode(email)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		emailVerificationCode = verificationCode.VerificationCode
	}

	order, err := parser.DecodeValid[*orders.OrderDTO](r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	order.UserId = userId
	order.Status = orders.StatusCreated

	orderId, err := h.ordersService.Create(order)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = h.emailAPI.Write(emailContent(
		r.Host,
		orderId,
		order.ExtraInformation,
		email,
		emailVerificationCode,
	))

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (h *HTTPHandler) getOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))
	if len(sessionKey) != sessionKeyLength {
		return http.StatusUnauthorized, errors.New("invalid session key")
	}

	user, sessionKey, err := h.authAPI.FindBySessionKey(sessionKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	order, err := h.ordersService.FindById(orderId, user.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, order, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	return 0, nil
}

func (h *HTTPHandler) verifyOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	if verificationCode := r.URL.Query().Get(verificationCodeQueryName); verificationCode != "" {
		email := r.URL.Query().Get(emailQueryName)

		sessionKey, err := h.authAPI.VerifyEmail(email, verificationCode)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))
	}

	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.ordersService.MarkAsVerified(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (h *HTTPHandler) completeOrder(_ http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.ordersService.MarkAsCompleted(id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (h *HTTPHandler) deleteOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))
	if len(sessionKey) != sessionKeyLength {
		return http.StatusUnauthorized, errors.New("invalid session key")
	}

	user, sessionKey, err := h.authAPI.FindBySessionKey(sessionKey)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.ordersService.Delete(orderId, user.Id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	return 0, nil
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

func emailContent(
	host string,
	orderId int,
	extraOrderInformation string,
	email string,
	verificationCode string,
) *emailapi.EmailDTO {
	var params string
	if verificationCode != "" {
		params = fmt.Sprintf("?%s=%s&%s=%s", emailQueryName, email, verificationCodeQueryName, verificationCode)
	}

	var url = fmt.Sprintf("http://%s/api/order/%d/verify%s", host, orderId, params)

	return &emailapi.EmailDTO{
		To:      email,
		Subject: "Order verification",
		Body: fmt.Sprintf(
			"\n\nYou've just made an order with following information:\n\n%s\n\nIf you didn't order this, delete this this email.\nFollow the link to verify your order. Don't share your link\n\n%s",
			extraOrderInformation, url,
		),
	}
}
