package messagesrepo

import "time"

type MessageEntity struct {
	Id             int
	Message        string
	SenderId       int
	ReadBySender   bool
	ReceiverId     int
	ReadByReceiver bool
	CreatedAt      time.Time
}
