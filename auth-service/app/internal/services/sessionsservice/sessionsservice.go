package sessionsservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/sessionsrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/str"
)

const (
	sessionKeyLength = 32
)

var (
	ErrSessionExpired  = errors.New("session is expired")
	ErrSessionNotFound = errors.New("session not found")
)

type SessionsService interface {
	Create(userId int) (string, error)
	CreateWithTTL(userId int, ttl time.Duration) (string, error)

	UpdateKey(sessionKey string) (string, error)

	Valid(sessionKey string) error
}

type sessionsService struct {
	ctx        context.Context
	defaultTTL time.Duration

	sessionsRepository sessionsrepo.SessionsRepository
}

func NewSessionsService(
	ctx context.Context,
	defaultTTL time.Duration,
	sessionsRepository sessionsrepo.SessionsRepository,
) SessionsService {
	return &sessionsService{
		ctx:                ctx,
		defaultTTL:         defaultTTL,
		sessionsRepository: sessionsRepository,
	}
}

func (s *sessionsService) Create(userId int) (string, error) {
	return s.CreateWithTTL(userId, s.defaultTTL)
}

func (s *sessionsService) CreateWithTTL(userId int, ttl time.Duration) (string, error) {
	session := str.Generate(sessionKeyLength)

	sessionEntity := &sessionsrepo.SessionEntity{
		UserId:     userId,
		SessionKey: session,
		ExpiresAt:  sql.NullTime{Time: time.Now().Add(ttl)},
	}

	err := s.sessionsRepository.Create(sessionEntity)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (s *sessionsService) UpdateKey(sessionKey string) (string, error) {
	newKey := str.Generate(sessionKeyLength)

	err := s.sessionsRepository.UpdateKey(sessionKey, newKey)
	if err != nil {
		return "", err
	}

	return newKey, nil
}

func (s *sessionsService) Valid(sessionKey string) error {
	entity, err := s.sessionsRepository.FindByKey(sessionKey)
	if errors.Is(err, sessionsrepo.ErrSessionNotFound) {
		return ErrSessionNotFound
	}
	if err != nil {
		return err
	}

	if time.Now().After(entity.ExpiresAt.Time) {
		return ErrSessionExpired
	}

	return nil
}

func (s *sessionsService) delete(id int) error {
	return s.sessionsRepository.Delete(id)
}
