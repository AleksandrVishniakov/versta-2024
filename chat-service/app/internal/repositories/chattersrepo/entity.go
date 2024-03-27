package chattersrepo

import "database/sql"

type ChatterEntity struct {
	Id          int
	UserId      sql.NullInt32
	TempSession sql.NullString
}

type ChatterEntityWithUnread struct {
	ChatterEntity
	UnreadMessagesCount int
}
