package messages

import (
	"time"

	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/messagesrepo"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/scrambler"
)

type MessageRequestDTO struct {
	Message        string `json:"message"`
	SenderId       int    `json:"senderId,omitempty"`
	ReceiverId     int    `json:"receiverId,omitempty"`
	ReadByReceiver bool   `json:"readByReceiver,omitempty"`
}

func MapEntityFromUserRequest(
	encryptor scrambler.Encryptor,
	dto *MessageRequestDTO,
) (*messagesrepo.MessageEntity, error) {
	msg, err := encryptor.Encrypt([]byte(dto.Message))
	if err != nil {
		return nil, err
	}

	return &messagesrepo.MessageEntity{
		Message:        string(msg),
		SenderId:       dto.SenderId,
		ReceiverId:     dto.ReceiverId,
		ReadByReceiver: dto.ReadByReceiver,
		CreatedAt:      time.Now(),
	}, nil
}

type MessageResponseDTO struct {
	Id             int       `json:"id"`
	Message        string    `json:"message"`
	SenderId       int       `json:"senderId"`
	ReadBySender   bool      `json:"readBySender"`
	ReceiverId     int       `json:"receiverId"`
	ReadByReceiver bool      `json:"readByReceiver"`
	CreatedAt      time.Time `json:"createdAt"`
}

func MapResponseFromEntity(
	decryptor scrambler.Decryptor,
	entity *messagesrepo.MessageEntity,
) (*MessageResponseDTO, error) {
	msg, err := decryptor.Decrypt([]byte(entity.Message))
	if err != nil {
		return nil, err
	}

	return &MessageResponseDTO{
		Id:             entity.Id,
		Message:        string(msg),
		SenderId:       entity.SenderId,
		ReadBySender:   entity.ReadBySender,
		ReceiverId:     entity.ReceiverId,
		ReadByReceiver: entity.ReadByReceiver,
		CreatedAt:      entity.CreatedAt,
	}, nil
}
