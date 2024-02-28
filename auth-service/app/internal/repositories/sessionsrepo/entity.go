package sessionsrepo

import (
	"database/sql"
	"time"
)

type SessionEntity struct {
	Id         int
	UserId     int
	SessionKey string
	CreatedAt  time.Time
	ExpiresAt  sql.NullTime
}
