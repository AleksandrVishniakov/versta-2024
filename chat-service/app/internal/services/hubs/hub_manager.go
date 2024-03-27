package hubs

import (
	"context"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/messages"
	"log"
	"strconv"
	"sync"
)

type HubManager struct {
	ctx  context.Context
	rw   sync.RWMutex
	hubs map[string]*Hub

	msgStorage messages.Storage
}

func NewHubManager(ctx context.Context, msgStorage messages.Storage) *HubManager {
	return &HubManager{
		ctx:        ctx,
		msgStorage: msgStorage,
		hubs:       make(map[string]*Hub),
	}
}

func (m *HubManager) GetAvailableHub(senderId, receiverId int) *Hub {
	m.rw.Lock()
	defer m.rw.Unlock()

	key := hubKey(senderId, receiverId)
	log.Println("hub request:", key)

	if hub, exists := m.hubs[key]; exists {
		return hub
	}

	hub := NewHub(m.msgStorage, senderId, receiverId)
	m.hubs[key] = hub

	go func() {
		ctx, cancel := context.WithCancel(m.ctx)
		defer cancel()
		err := hub.Run(ctx)
		log.Printf("hub %s finished work wis error: %s", key, err)

		m.rw.Lock()
		defer m.rw.Unlock()
		if _, ok := m.hubs[key]; ok {
			delete(m.hubs, key)
		}
	}()

	return hub
}

func hubKey(senderId, receiverId int) string {
	return strconv.Itoa(min(senderId, receiverId)) +
		":" +
		strconv.Itoa(max(senderId, receiverId))
}
