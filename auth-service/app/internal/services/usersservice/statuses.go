package usersservice

type UserStatus string

const (
	StatusUser  UserStatus = "user"
	StatusAdmin UserStatus = "admin"
)
