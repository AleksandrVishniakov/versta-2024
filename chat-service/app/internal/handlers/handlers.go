package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/chatters"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/chattokens"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/hubs"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/messages"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/jwttokens"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/parser"
	"github.com/gorilla/websocket"
)

const (
	receiverId = 1
)

type HTTPHandler struct {
	ctx              context.Context
	upgrader         *websocket.Upgrader
	hubManager       *hubs.HubManager
	jwtTokensManager jwttokens.Manager
	chatTokens       *chattokens.ChatTokens

	chattersStorage chatters.Storage
	messagesStorage messages.Storage
}

func NewHTTPHandler(
	ctx context.Context,
	hubManager *hubs.HubManager,
	jwtTokensManager jwttokens.Manager,
	chattersStorage chatters.Storage,
	messagesStorage messages.Storage,
	chatTokens *chattokens.ChatTokens,
) *HTTPHandler {
	return &HTTPHandler{
		ctx:              ctx,
		hubManager:       hubManager,
		jwtTokensManager: jwtTokensManager,
		chattersStorage:  chattersStorage,
		messagesStorage:  messagesStorage,
		chatTokens:       chatTokens,
		upgrader: &websocket.Upgrader{
			HandshakeTimeout: 1 * time.Second,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
			CheckOrigin: func(*http.Request) bool {
				return true
			},
		},
	}
}

func (h *HTTPHandler) Handler() http.Handler {
	var mux = http.NewServeMux()

	var auth = NewChatterAuth(
		h.jwtTokensManager,
		h.chattersStorage,
	)
	var adminFromToken = NewAdminOnlyFromQueryToken(
		h.jwtTokensManager,
	)

	mux.HandleFunc("GET /ping", h.pingHandler)

	mux.Handle("GET /api/chat/preflight", auth(Errors(h.getChatToken)))

	mux.Handle("GET /api/chat", Errors(h.connectChat))
	mux.Handle("GET /api/messages", auth(Errors(h.getAllMessages)))
	mux.Handle("GET /api/messages/unread", auth(Errors(h.getUnreadMessagesCount)))
	mux.Handle("GET /api/messages/read_all", auth(Errors(h.readAllMessages)))

	mux.Handle("GET /api/admin/chat", adminFromToken(Errors(h.connectAdminChat)))
	mux.Handle("GET /api/admin/clients", auth(AdminOnly(Errors(h.getAdminClients))))
	mux.Handle("GET /api/admin/messages", auth(AdminOnly(Errors(h.getAllAdminMessages))))
	mux.Handle("GET /api/admin/messages/unread", auth(AdminOnly(Errors(h.getAdminUnreadMessagesCount))))
	mux.Handle("GET /api/admin/messages/read_all", auth(AdminOnly(Errors(h.readAllAdminMessages))))

	return Recovery(Logger(CORS(mux)))
}

func (h *HTTPHandler) pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) connectChat(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatToken := r.URL.Query().Get("t")
	if chatToken == "" {
		return http.StatusBadRequest, errors.New("empty chat token parameter")
	}

	chatterId, err := h.chatTokens.Get(chatToken)
	if errors.Is(err, chattokens.ErrExpiredToken) || errors.Is(err, chattokens.ErrNonExistentKey) {
		return http.StatusPreconditionRequired, err
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws connection failed:", err.Error())
		return http.StatusUpgradeRequired, err
	}

	ctx, cancel := context.WithCancel(h.ctx)
	hub := h.hubManager.GetAvailableHub(chatterId, receiverId)
	client := hubs.NewWSClient(hub, conn, chatterId)

	go func() {
		defer cancel()
		err := client.Listen(ctx)
		log.Println("client stop with error:", err)
	}()

	return http.StatusOK, nil
}

func (h *HTTPHandler) getAllMessages(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	msgs, err := h.messagesStorage.FindByChatterId(chatterId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, msgs, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) getChatToken(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	chatToken, err := h.chatTokens.Create(chatterId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	type ChatTokenResponse struct {
		ChatterId int    `json:"chatterId"`
		ChatToken string `json:"chatToken"`
	}

	err = parser.EncodeResponse(w, &ChatTokenResponse{
		ChatterId: chatterId,
		ChatToken: chatToken,
	}, http.StatusOK)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) getUnreadMessagesCount(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	count, err := h.messagesStorage.GetUnreadCount(chatterId, receiverId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, count, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) getAdminClients(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	clients, err := h.chattersStorage.FindSendersByChatterId(chatterId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, clients, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) connectAdminChat(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatToken := r.URL.Query().Get("t")
	if chatToken == "" {
		return http.StatusBadRequest, errors.New("empty chat token parameter")
	}

	withId, err := strconv.Atoi(r.URL.Query().Get("with"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	chatterId, err := h.chatTokens.Get(chatToken)
	if errors.Is(err, chattokens.ErrExpiredToken) || errors.Is(err, chattokens.ErrNonExistentKey) {
		return http.StatusPreconditionRequired, err
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws connection failed:", err.Error())
		return http.StatusUpgradeRequired, err
	}

	ctx, cancel := context.WithCancel(h.ctx)
	hub := h.hubManager.GetAvailableHub(chatterId, withId)
	client := hubs.NewWSClient(hub, conn, chatterId)

	go func() {
		defer cancel()
		err := client.Listen(ctx)
		log.Println("client stop with error:", err)
	}()

	return http.StatusOK, nil
}

func (h *HTTPHandler) getAllAdminMessages(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	withId, err := strconv.Atoi(r.URL.Query().Get("with"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	msgs, err := h.messagesStorage.FindBySenderAndReceiver(chatterId, withId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, msgs, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) getAdminUnreadMessagesCount(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	withId, err := strconv.Atoi(r.URL.Query().Get("with"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	count, err := h.messagesStorage.GetUnreadCount(chatterId, withId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = parser.EncodeResponse(w, count, http.StatusOK)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) readAllMessages(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.messagesStorage.ReadAll(chatterId, receiverId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (h *HTTPHandler) readAllAdminMessages(w http.ResponseWriter, r *http.Request) (statusCode int, err error) {
	withId, err := strconv.Atoi(r.URL.Query().Get("with"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	chatterId, err := strconv.Atoi(fmt.Sprintf("%v", r.Context().Value(ChatterIdContextKey)))
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = h.messagesStorage.ReadAll(chatterId, withId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
