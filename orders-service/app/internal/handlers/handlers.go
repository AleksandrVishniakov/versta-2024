package handlers

import (
	"net/http"
)

func NewHTTPHandler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)

	mux.Handle("GET /api/orders", Errors(getAllOrders))
	mux.Handle("POST /api/order", Errors(createNewOrder))
	mux.Handle("GET /api/order/{orderId}", Errors(getOrder))
	mux.Handle("DELETE /api/order/{orderId}", Errors(deleteOrder))
	mux.Handle("PUT /api/order/{orderId}/verify", Errors(verifyOrder))
	mux.Handle("PUT /api/order/{orderId}/complete", Errors(completeOrder))

	return Recovery(Logger(mux))
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
