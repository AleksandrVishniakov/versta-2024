package hubs

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/messages"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/types"
)

type Hub struct {
	msgStorage messages.Storage

	rw      sync.RWMutex
	clients map[string]*WSClient

	chattersId *types.Pair[int, int]

	broadcast  chan *messages.MessageRequestDTO
	register   chan *WSClient
	unregister chan *WSClient
}

func NewHub(
	msgStorage messages.Storage,
	chatterId1, chatterId2 int,
) *Hub {
	return &Hub{
		msgStorage: msgStorage,
		clients:    make(map[string]*WSClient),
		chattersId: &types.Pair[int, int]{First: chatterId1, Second: chatterId2},
		broadcast:  make(chan *messages.MessageRequestDTO),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

func (h *Hub) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
			if len(h.clients) == 0 {
				return nil
			}
		case message := <-h.broadcast:
			message.ReceiverId = h.detectMsgReceiver(message.SenderId)

			slog.Debug("message",
				slog.Int("senderId", message.SenderId),
				slog.Int("receiverId", message.ReceiverId),
				slog.String("message", message.Message),
			)

			err := h.notifyClients(message)
			if err != nil {
				return err
			}
		}
	}
}

func (h *Hub) registerClient(c *WSClient) {
	h.rw.Lock()
	defer h.rw.Unlock()

	log.Printf("new chatter #%d with uuid=%s", c.ChatterId, c.UUID)

	h.clients[c.UUID] = c
}

func (h *Hub) unregisterClient(c *WSClient) {
	h.rw.Lock()
	defer h.rw.Unlock()

	if _, ok := h.clients[c.UUID]; ok {
		delete(h.clients, c.UUID)
		log.Printf("delete chatter #%d with uuid=%s", c.ChatterId, c.UUID)
	}
}

func (h *Hub) notifyClients(msg *messages.MessageRequestDTO) error {
	id, err := h.msgStorage.Create(
		msg.Message,
		msg.SenderId,
		msg.ReceiverId,
		h.isChatterOnline(msg.ReceiverId),
	)
	if err != nil {
		return err
	}

	message := &messages.MessageResponseDTO{
		Id:         id,
		Message:    msg.Message,
		SenderId:   msg.SenderId,
		ReceiverId: msg.ReceiverId,
		CreatedAt:  time.Now(),
	}

	h.rw.Lock()
	defer h.rw.Unlock()

	for _, c := range h.clients {
		c.send <- message
	}

	return nil
}

func (h *Hub) isChatterOnline(id int) bool {
	h.rw.RLock()
	defer h.rw.RUnlock()

	for _, c := range h.clients {
		if c.ChatterId == id {
			return true
		}
	}

	return false
}

func (h *Hub) detectMsgReceiver(senderId int) int {
	switch senderId {
	case h.chattersId.First:
		return h.chattersId.Second
	case h.chattersId.Second:
		return h.chattersId.First
	default:
		return 0
	}
}
