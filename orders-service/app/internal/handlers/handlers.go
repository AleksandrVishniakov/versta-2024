package handlers

import (
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/jwttokens"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/internal/services/orders"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/parser"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/utils"
)

const (
	emailQueryName            = "email"
	verificationCodeQueryName = "code"
)

type HTTPHandler struct {
	ordersService   orders.Service
	tokensManager   jwttokens.Manager
	refreshTokenTTL time.Duration
}

func NewHTTPHandler(
	ordersService orders.Service,
	tokensManager jwttokens.Manager,
	refreshTokenTTL time.Duration,
) *HTTPHandler {
	return &HTTPHandler{
		ordersService:   ordersService,
		tokensManager:   tokensManager,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	jwtAuth := NewJWTAuthMiddleware(h.tokensManager)

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.Handle("GET /api/orders", jwtAuth(Errors(h.getAllOrders)))
	mux.Handle("POST /api/order", Errors(h.createNewOrder))
	mux.Handle("GET /api/order/{orderId}", jwtAuth(Errors(h.getOrder)))
	mux.Handle("DELETE /api/order/{orderId}", jwtAuth(Errors(h.deleteOrder)))
	mux.Handle("GET /api/order/{orderId}/verify", Errors(h.verifyOrder))

	mux.Handle("GET /api/order/{orderId}/complete", Errors(h.completeOrder))

	return Recovery(Logger(CORS(mux)))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) getAllOrders(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	email := fmt.Sprintf("%v", r.Context().Value(EmailContextKey))

	userOrders, err := h.ordersService.FindAll(email)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	err = parser.EncodeResponse(w, userOrders, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return 0, nil
}

func (h *HTTPHandler) createNewOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

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

	orderId, err := h.ordersService.Create(email, orderReq.ExtraInformation)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	err = parser.EncodeResponse(w, orderId, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return 0, nil
}

func (h *HTTPHandler) getOrder(w http.ResponseWriter, r *http.Request) (int, error) {
	defer utils.CloseReadCloser(r.Body)

	email := fmt.Sprintf("%v", r.Context().Value(EmailContextKey))

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	order, err := h.ordersService.FindById(email, orderId)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

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

	accessToken, refreshToken, err := h.ordersService.Verify(email, orderId, code)
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

	if refreshToken != "" {
		http.SetCookie(w, refreshTokenCookie(refreshToken, h.refreshTokenTTL))
	}

	if accessToken != "" {
		err := parser.EncodeResponse(w, accessToken, http.StatusOK)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

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

	email := fmt.Sprintf("%v", r.Context().Value(EmailContextKey))

	orderId, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.ordersService.Delete(email, orderId)
	if errors.Is(err, orders.ErrWithAuthorization) {
		return http.StatusUnauthorized, err
	}
	if errors.Is(err, orders.ErrNoOrders) {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return 0, nil
}

func refreshTokenCookie(token string, ttl time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshTokenCookieKey,
		Value:    token,
		Expires:  time.Now().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
		HttpOnly: true,
	}
}
