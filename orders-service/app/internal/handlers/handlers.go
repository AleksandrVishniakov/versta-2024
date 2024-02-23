package handlers

import (
	"net/http"
)

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) InitRoutes() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", Logger(pingHandler))

	mux.HandleFunc("GET /api/orders", Recovery(Logger(Errors(getAllOrders))))
	mux.HandleFunc("POST /api/order", Recovery(Logger(Errors(createNewOrder))))
	mux.HandleFunc("GET /api/order/{orderId}", Recovery(Logger(Errors(getOrder))))
	mux.HandleFunc("DELETE /api/order/{orderId}", Recovery(Logger(Errors(deleteOrder))))
	mux.HandleFunc("PUT /api/order/{orderId}/verify", Recovery(Logger(Errors(verifyOrder))))
	mux.HandleFunc("PUT /api/order/{orderId}/complete", Recovery(Logger(Errors(completeOrder))))

	return mux
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getAllOrders(w http.ResponseWriter, r *http.Request) error {
	// TODO implement getAllOrders
	panic("not implemented")
}

func createNewOrder(w http.ResponseWriter, r *http.Request) error {
	// TODO implement createNewOrder
	panic("not implemented")
}

func getOrder(w http.ResponseWriter, r *http.Request) error {
	// TODO implement getOrder
	panic("not implemented")
}

func verifyOrder(w http.ResponseWriter, r *http.Request) error {
	// TODO implement verifyOrder
	panic("not implemented")
}

func completeOrder(w http.ResponseWriter, r *http.Request) error {
	// TODO implement completeOrder
	panic("not implemented")
}

func deleteOrder(w http.ResponseWriter, r *http.Request) error {
	// TODO implement deleteOrder
	panic("not implemented")
}
