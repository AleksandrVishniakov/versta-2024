package usersrepo

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	Id                    int
	Email                 string
	Name                  sql.NullString
	IsEmailVerified       bool
	EmailVerificationCode sql.NullString
	CreatedAt             time.Time
}
