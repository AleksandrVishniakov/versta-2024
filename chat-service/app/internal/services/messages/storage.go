package messages

import (
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/messagesrepo"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/scrambler"
)

type Storage interface {
	Create(message string, senderId, receiverId int, readByReceiver bool) (int, error)

	FindByChatterId(chatterId int) ([]*MessageResponseDTO, error)
	FindBySenderAndReceiver(senderId, receiverId int) ([]*MessageResponseDTO, error)

	GetUnreadCount(forId, withId int) (int, error)

	ReadAll(forId, withId int) error
}

type storage struct {
	repository messagesrepo.MessagesRepository
	scrambler  scrambler.Scrambler
}

func NewMessagesStorage(
	repository messagesrepo.MessagesRepository,
	scrambler scrambler.Scrambler,
) Storage {
	return &storage{
		repository: repository,
		scrambler:  scrambler,
	}
}

func (s *storage) Create(message string, senderId, receiverId int, readByReceiver bool) (int, error) {
	entity, err := MapEntityFromUserRequest(s.scrambler, &MessageRequestDTO{
		Message:        message,
		SenderId:       senderId,
		ReceiverId:     receiverId,
		ReadByReceiver: readByReceiver,
	})

	if err != nil {
		return 0, err
	}

	id, err := s.repository.Create(entity)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *storage) FindByChatterId(chatterId int) (messages []*MessageResponseDTO, err error) {
	entities, err := s.repository.FindByChatterId(chatterId)
	if err != nil {
		return nil, err
	}

	messages = []*MessageResponseDTO{}

	messages, err = s.parseEntities(entities, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *storage) FindBySenderAndReceiver(senderId, receiverId int) (messages []*MessageResponseDTO, err error) {
	entities, err := s.repository.FindBySenderAndReceiver(senderId, receiverId)
	if err != nil {
		return nil, err
	}

	messages = []*MessageResponseDTO{}

	messages, err = s.parseEntities(entities, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *storage) GetUnreadCount(forId, withId int) (int, error) {
	return s.repository.GetUnreadCount(forId, withId)
}

func (s *storage) ReadAll(forId, withId int) error {
	return s.repository.ReadAll(forId, withId)
}

func (s *storage) parseEntities(entities []*messagesrepo.MessageEntity, messages []*MessageResponseDTO) ([]*MessageResponseDTO, error) {
	for _, e := range entities {
		dto, err := MapResponseFromEntity(s.scrambler, e)
		if err != nil {
			return nil, err
		}

		messages = append(messages, dto)
	}
	return messages, nil
}
