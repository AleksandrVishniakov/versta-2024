package handlers

import (
	"net/http"
)

func NewHTTPHandler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)

	mux.HandleFunc("GET /api/orders", getAllOrders)
	mux.HandleFunc("POST /api/order", createNewOrder)
	mux.HandleFunc("GET /api/order/{orderId}", getOrder)
	mux.HandleFunc("DELETE /api/order/{orderId}", deleteOrder)
	mux.HandleFunc("PUT /api/order/{orderId}/verify", verifyOrder)
	mux.HandleFunc("PUT /api/order/{orderId}/complete", completeOrder)

	return Recovery(Logger(mux))
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	// TODO implement getAllOrders
	panic("not implemented")
}

func createNewOrder(w http.ResponseWriter, r *http.Request) {
	// TODO implement createNewOrder
	panic("not implemented")
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	// TODO implement getOrder
	panic("not implemented")
}

func verifyOrder(w http.ResponseWriter, r *http.Request) {
	// TODO implement verifyOrder
	panic("not implemented")
}

func completeOrder(w http.ResponseWriter, r *http.Request) {
	// TODO implement completeOrder
	panic("not implemented")
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	// TODO implement deleteOrder
	panic("not implemented")
}
