package ordersrepo

import "database/sql"

type OrderEntity struct {
	Id               int
	UserId           int
	ExtraInformation sql.NullString
	Status           byte
	VerificationCode sql.NullString
}
