package hubs

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/services/messages"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/str"
	"github.com/gorilla/websocket"
)

const (
	writeTimeout = 10 * time.Second
	pongTime     = 5 * time.Second
	pingTime     = (pongTime * 9) / 10
)

type WSClient struct {
	ChatterId int
	UUID      string

	hub *Hub

	mu         *sync.Mutex
	connection *websocket.Conn

	send chan *messages.MessageResponseDTO
}

func NewWSClient(hub *Hub, conn *websocket.Conn, id int) *WSClient {
	client := &WSClient{
		ChatterId:  id,
		UUID:       str.Generate(64),
		hub:        hub,
		mu:         &sync.Mutex{},
		connection: conn,
		send:       make(chan *messages.MessageResponseDTO),
	}

	client.register()

	return client
}

func (c *WSClient) Listen(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	readErrChan := c.readPump(ctx)
	writeErrChan := c.writePump(ctx)

	defer cancel()
	defer func() {
		c.unregister()
		close(c.send)
		if err := c.connection.Close(); err != nil {
			slog.Error("connection writing error",
				slog.String("error", err.Error()),
			)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readErrChan:
			if err != nil {
				return e.WrapErr(err, "wsclient reading")
			}
		case err := <-writeErrChan:
			if err != nil {
				return e.WrapErr(err, "wsclient writing")
			}
		}
	}
}

func (c *WSClient) register() {
	c.hub.register <- c
}

func (c *WSClient) unregister() {
	c.hub.unregister <- c
}

func (c *WSClient) readPump(ctx context.Context) <-chan error {
	errChan := make(chan error)

	if err := c.connection.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
		errChan <- e.WrapErr(err, "set read deadline error")
	}
	c.connection.SetPongHandler(func(string) error {
		if err := c.connection.SetReadDeadline(time.Now().Add(pongTime)); err != nil {
			errChan <- e.WrapErr(err, "set read deadline error")
		}
		return nil
	})

	go func() {
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				mt, message, err := c.connection.ReadMessage()
				if mt == websocket.CloseMessage {
					break
				}

				if err != nil {
					errChan <- e.WrapErr(err, "message reading error")
					continue
				}

				message = bytes.TrimSpace(message)

				msg := &messages.MessageRequestDTO{
					Message:  string(message),
					SenderId: c.ChatterId,
				}

				c.hub.broadcast <- msg
			}
		}
	}()

	return errChan
}

func (c *WSClient) writePump(ctx context.Context) <-chan error {
	var errChan = make(chan error)

	go func() {
		defer close(errChan)
		ticker := time.NewTicker(pingTime)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case message, ok := <-c.send:
				if !ok {
					return
				}

				c.mu.Lock()
				err := c.connection.WriteJSON(message)
				c.mu.Unlock()

				if err != nil {
					errChan <- e.WrapErr(err, "message writing error")
				}
			case <-ticker.C:
				c.mu.Lock()
				err := c.write(websocket.PingMessage, []byte{})
				c.mu.Unlock()
				if err != nil {
					errChan <- e.WrapErr(err, "ping-message writing error")
				}
			}
		}
	}()

	return errChan
}

func (c *WSClient) write(mt int, msg []byte) error {
	if err := c.connection.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}

	return c.connection.WriteMessage(mt, msg)
}
