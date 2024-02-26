package orders

type OrderStatus byte

const (
	StatusCreated OrderStatus = iota
	StatusVerified
	StatusCompleted
)
