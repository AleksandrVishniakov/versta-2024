package chatters

import (
	"database/sql"
	"errors"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/chattersrepo"
)

var (
	ErrChatterNotFound = errors.New("chatters: chatter not found")
)

type Storage interface {
	CreateWithId(userId int) (int, error)
	CreateWithSession(session string) (int, error)

	FindBySession(session string) (*ChatterDTO, error)
	FindByChatterId(chatterId int) (*ChatterDTO, error)
	FindByUserId(userId int) (*ChatterDTO, error)
	FindSendersByChatterId(chatterId int) ([]*ChatterWithUnreadDTO, error)

	ChangeSessionToId(session string, userId int) error
}

type storage struct {
	repository chattersrepo.ChattersRepository
}

func NewChattersStorage(repository chattersrepo.ChattersRepository) Storage {
	return &storage{repository: repository}
}

func (s *storage) CreateWithId(userId int) (int, error) {
	return s.repository.Create(&chattersrepo.ChatterEntity{
		UserId: sql.NullInt32{Int32: int32(userId)},
	})
}

func (s *storage) CreateWithSession(session string) (int, error) {
	return s.repository.Create(&chattersrepo.ChatterEntity{
		TempSession: sql.NullString{String: session},
	})
}

func (s *storage) FindBySession(session string) (*ChatterDTO, error) {
	entity, err := s.repository.FindBySession(session)
	if errors.Is(err, chattersrepo.ErrChatterNotFound) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return MapChatterDTOFromEntity(entity), nil
}

func (s *storage) FindByChatterId(chatterId int) (*ChatterDTO, error) {
	entity, err := s.repository.FindByChatterId(chatterId)
	if errors.Is(err, chattersrepo.ErrChatterNotFound) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return MapChatterDTOFromEntity(entity), nil
}

func (s *storage) FindByUserId(userId int) (*ChatterDTO, error) {
	entity, err := s.repository.FindByUserId(userId)
	if errors.Is(err, chattersrepo.ErrChatterNotFound) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return MapChatterDTOFromEntity(entity), nil
}

func (s *storage) ChangeSessionToId(session string, userId int) error {
	return s.repository.ChangeSessionToId(session, userId)
}

func (s *storage) FindSendersByChatterId(chatterId int) (chatters []*ChatterWithUnreadDTO, err error) {
	entities, err := s.repository.FindSendersByChatterId(chatterId)
	if err != nil {
		return nil, err
	}

	chatters = []*ChatterWithUnreadDTO{}

	for _, e := range entities {
		chatters = append(chatters, MapChatterWithUnreadDTOFromEntity(e))
	}

	return chatters, nil
}
