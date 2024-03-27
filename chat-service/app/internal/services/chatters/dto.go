package chatters

import "github.com/AleksandrVishniakov/versta-2024/chat-service/app/internal/repositories/chattersrepo"

type ChatterDTO struct {
	Id          int    `json:"id"`
	UserId      int    `json:"userId"`
	TempSession string `json:"tempSession"`
}

func MapChatterDTOFromEntity(entity *chattersrepo.ChatterEntity) *ChatterDTO {
	return &ChatterDTO{
		Id:          entity.Id,
		UserId:      int(entity.UserId.Int32),
		TempSession: entity.TempSession.String,
	}
}

type ChatterWithUnreadDTO struct {
	Id                  int    `json:"id"`
	UserId              int    `json:"userId"`
	TempSession         string `json:"tempSession"`
	UnreadMessagesCount int    `json:"unreadMessagesCount"`
}

func MapChatterWithUnreadDTOFromEntity(entity *chattersrepo.ChatterEntityWithUnread) *ChatterWithUnreadDTO {
	return &ChatterWithUnreadDTO{
		Id:                  entity.Id,
		UserId:              int(entity.UserId.Int32),
		TempSession:         entity.TempSession.String,
		UnreadMessagesCount: entity.UnreadMessagesCount,
	}
}
