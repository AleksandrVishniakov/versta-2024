package usersservice

import (
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/encryptor"
	"time"
)

type UserResponseDTO struct {
	Id              int       `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	CreatedAt       time.Time `json:"createdAt"`
}

type VerificationCodeDTO struct {
	VerificationCode string `json:"verificationCode"`
}

func mapResponseFromEntity(crypt *encryptor.Encryptor, entity *usersrepo.UserEntity) (*UserResponseDTO, error) {
	email, err := crypt.Decrypt([]byte(entity.Email))
	if err != nil {
		return nil, err
	}

	var name string
	if entity.Name.Valid {
		nameBytes, err := crypt.Decrypt([]byte(entity.Name.String))
		if err != nil {
			return nil, err
		}

		name = string(nameBytes)
	}

	return &UserResponseDTO{
		Id:              entity.Id,
		Email:           string(email),
		Name:            name,
		IsEmailVerified: entity.IsEmailVerified,
		CreatedAt:       entity.CreatedAt,
	}, nil
}
