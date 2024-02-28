package auth

type UserDTO struct {
	Id              int    `json:"id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	IsEmailVerified bool   `json:"isEmailVerified"`
}
