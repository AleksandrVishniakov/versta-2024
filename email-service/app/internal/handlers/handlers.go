package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/AleksandrVishniakov/versta-2024/email-service/app/internal/services/emailservice"
	"github.com/AleksandrVishniakov/versta-2024/email-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/email-service/app/pkg/parser"
)

type HTTPHandler struct {
	emailService emailservice.EmailService
}

func NewHTTPHandler(emailService emailservice.EmailService) *HTTPHandler {
	return &HTTPHandler{emailService: emailService}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	mux.HandleFunc("/ping", h.pingHandler)

	mux.HandleFunc("POST /api/email", h.handleSendEmailRequest)

	return Recovery(Logger(mux))
}

func (h *HTTPHandler) handleSendEmailRequest(w http.ResponseWriter, r *http.Request) {
	emailContent, err := parser.DecodeValid[emailservice.EmailDTO](r.Body)
	if err != nil {
		e.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer closeReadCloser(r.Body)

	err = h.emailService.Write(&emailContent)
	if err != nil {
		e.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
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
