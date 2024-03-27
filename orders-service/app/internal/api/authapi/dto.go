package authapi

type UserDTO struct {
	Id              int        `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Status          UserStatus `json:"status"`
	IsEmailVerified bool       `json:"isEmailVerified"`
}

type VerificationCodeDTO struct {
	VerificationCode string `json:"verificationCode"`
}

type UserStatus string

const (
	StatusUser  UserStatus = "user"
	StatusAdmin UserStatus = "admin"
)
