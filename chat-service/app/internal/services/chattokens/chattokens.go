package chattokens

import (
	"errors"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/str"
)

const (
	keyLength = 16
	valueTTL  = 15 * time.Minute
)

var (
	ErrExpiredToken   = errors.New("chattokens: token is expired")
	ErrNonExistentKey = errors.New("chattokens: key is non-existent")
)

type tokenValue struct {
	chatterId int
	expiresAt time.Time
}

type ChatTokens struct {
	mutex   sync.RWMutex
	storage map[string]*tokenValue
}

func NewChatTokens() *ChatTokens {
	return &ChatTokens{
		mutex:   sync.RWMutex{},
		storage: make(map[string]*tokenValue),
	}
}

func (c *ChatTokens) Create(chatterId int) (string, error) {
	key := str.Generate(keyLength)

	value := &tokenValue{
		chatterId: chatterId,
		expiresAt: time.Now().Add(valueTTL),
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.storage[key]; exists {
		return "", errors.New("chattokens: key already exists")
	}

	c.storage[key] = value

	return key, nil
}

func (c *ChatTokens) Get(key string) (int, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok := c.storage[key]
	if !ok {
		return 0, ErrNonExistentKey
	}

	if time.Now().After(value.expiresAt) {
		delete(c.storage, key)
		return 0, ErrExpiredToken
	}

	return value.chatterId, nil
}
