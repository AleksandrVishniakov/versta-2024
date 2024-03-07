package orders

import (
	"database/sql"
	"fmt"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/repositories/ordersrepo"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/scrambler"
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

func MapEntityFromDTO(order *OrderDTO, encryptor scrambler.Encryptor) (*ordersrepo.OrderEntity, error) {
	encryptedInfo, err := encryptor.Encrypt([]byte(order.ExtraInformation))
	if err != nil {
		return nil, err
	}

	return &ordersrepo.OrderEntity{
		Id:               order.Id,
		UserId:           order.UserId,
		ExtraInformation: sql.NullString{String: string(encryptedInfo)},
		Status:           byte(order.Status),
	}, nil
}

func MapDTOFromEntity(entity *ordersrepo.OrderEntity, decryptor scrambler.Decryptor) (*OrderDTO, error) {
	decryptedInfo, err := decryptor.Decrypt([]byte(entity.ExtraInformation.String))
	if err != nil {
		return nil, err
	}

	return &OrderDTO{
		Id:               entity.Id,
		UserId:           entity.UserId,
		ExtraInformation: string(decryptedInfo),
		Status:           OrderStatus(entity.Status),
	}, nil
}
