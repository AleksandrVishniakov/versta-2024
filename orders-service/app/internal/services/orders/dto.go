package orders

import (
	"database/sql"
	"fmt"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/repositories/ordersrepo"
)

type OrderDTO struct {
	Id               int         `json:"id"`
	UserId           int         `json:"userId"`
	ExtraInformation string      `json:"extraInformation"`
	Status           OrderStatus `json:"status"`
}

func (o *OrderDTO) Valid() (bool, error) {
	if o.Status != StatusCreated {
		return false, fmt.Errorf("unknown status %d", o.Status)
	}

	return true, nil
}

func MapEntityFromDTO(order *OrderDTO) *ordersrepo.OrderEntity {
	return &ordersrepo.OrderEntity{
		Id:               order.Id,
		UserId:           order.UserId,
		ExtraInformation: sql.NullString{String: order.ExtraInformation},
		Status:           byte(order.Status),
	}
}

func MapDTOFromEntity(entity *ordersrepo.OrderEntity) *OrderDTO {
	return &OrderDTO{
		Id:               entity.Id,
		UserId:           entity.UserId,
		ExtraInformation: entity.ExtraInformation.String,
		Status:           OrderStatus(entity.Status),
	}
}
