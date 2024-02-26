package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/orders"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/ordersservice"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/e"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/parser"
)

const userId = 1 // TODO: replace const with data from cookies and email auth service

type HTTPHandler struct {
	ordersService ordersservice.OrdersService
}

func NewHTTPHandler(
	ordersService ordersservice.OrdersService,
) *HTTPHandler {
	return &HTTPHandler{
		ordersService: ordersService,
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.HandleFunc("GET /api/orders", h.getAllOrders)
	mux.HandleFunc("POST /api/order", h.createNewOrder)
	mux.HandleFunc("GET /api/order/{orderId}", h.getOrder)
	mux.HandleFunc("DELETE /api/order/{orderId}", h.deleteOrder)
	mux.HandleFunc("PUT /api/order/{orderId}/verify", h.verifyOrder)
	mux.HandleFunc("PUT /api/order/{orderId}/complete", h.completeOrder)

	return Recovery(Logger(mux))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) getAllOrders(w http.ResponseWriter, _ *http.Request) {
	userOrders, err := h.ordersService.FindAll(userId)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = parser.EncodeResponse(w, userOrders, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) createNewOrder(w http.ResponseWriter, r *http.Request) {
	defer closeReadCloser(r.Body)

	order, err := parser.DecodeValid[*orders.OrderDTO](r.Body)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order.UserId = userId
	order.Status = orders.StatusCreated

	_, err = h.ordersService.Create(order)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.ordersService.FindById(id, userId)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = parser.EncodeResponse(w, order, http.StatusOK)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) verifyOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ordersService.MarkAsVerified(id, userId)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) completeOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ordersService.MarkAsCompleted(id, userId)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) deleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("orderId"))
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.ordersService.Delete(id, userId)
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
