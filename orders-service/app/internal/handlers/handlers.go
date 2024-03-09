package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/orders"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/parser"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/utils"
)

const (
	sessionKeyLength          = 32
	emailQueryName            = "email"
	verificationCodeQueryName = "code"
)

type HTTPHandler struct {
	ordersService orders.Service
	cookieTTL     time.Duration
}

func NewHTTPHandler(
	ordersService orders.Service,
	cookieTTL time.Duration,
) *HTTPHandler {
	return &HTTPHandler{
		ordersService: ordersService,
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

	return Recovery(Logger(CORS(Cookies(mux))))
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

	userOrders, sessionKey, err := h.ordersService.FindAll(sessionKey)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	err = parser.EncodeResponse(w, userOrders, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (h *HTTPHandler) createNewOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)
	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))

	email := r.URL.Query().Get(emailQueryName)
	if email == "" {
		return http.StatusBadRequest, errors.New("empty email")
	}

	type orderRequest struct {
		ExtraInformation string `json:"extraInformation"`
	}

	orderReq, err := parser.Decode[orderRequest](r.Body)
	if err != nil {
		return http.StatusBadRequest, err
	}

	orderId, sessionKey, err := h.ordersService.Create(sessionKey, email, orderReq.ExtraInformation)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	if sessionKey != "" {
		http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))
	}

	err = parser.EncodeResponse(w, orderId, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return 0, nil
}

func (h *HTTPHandler) getOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))
	if len(sessionKey) != sessionKeyLength {
		return http.StatusUnauthorized, errors.New("invalid session key")
	}

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	order, sessionKey, err := h.ordersService.FindById(sessionKey, orderId)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	err = parser.EncodeResponse(w, order, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return 0, nil
}

func (h *HTTPHandler) verifyOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	email := r.URL.Query().Get(emailQueryName)
	if email == "" {
		return http.StatusBadRequest, errors.New("empty email")
	}

	code := r.URL.Query().Get(verificationCodeQueryName)
	if len(code) != 6 {
		return http.StatusBadRequest, errors.New("invalid verification code length")
	}

	sessionKey, err := h.ordersService.Verify(email, orderId, code)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if errors.Is(err, orders.ErrMismatchedCodes) {
		return http.StatusBadRequest, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}

	http.SetCookie(w, sessionCookie(sessionKey, h.cookieTTL))

	return 0, nil
}

func (h *HTTPHandler) completeOrder(_ http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.ordersService.Complete(orderId)
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return 0, nil
}

func (h *HTTPHandler) deleteOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	sessionKey := fmt.Sprintf("%v", r.Context().Value(SessionKey))
	if len(sessionKey) != sessionKeyLength {
		return http.StatusUnauthorized, errors.New("invalid session key")
	}

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	sessionKey, err = h.ordersService.Delete(sessionKey, orderId)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
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
