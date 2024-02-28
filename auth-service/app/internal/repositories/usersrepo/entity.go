package usersrepo

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	Id              int
	Email           string
	Name            sql.NullString
	IsEmailVerified bool
	CreatedAt       time.Time
}
