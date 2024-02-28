package auth

type API interface {
	FindById(userId int) (*UserDTO, error)
}
