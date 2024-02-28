package auth

type Mock struct {
	data []*UserDTO
}

func NewMock(data []*UserDTO) API {
	if data == nil {
		data = []*UserDTO{
			{
				Id:              0,
				Email:           "example@example.com",
				Name:            "John",
				IsEmailVerified: true,
			},
			{
				Id:              1,
				Email:           "jane@gmail.com",
				Name:            "Jane",
				IsEmailVerified: true,
			},
			{
				Id:              2,
				Email:           "pete.jackson@gmail.com",
				Name:            "Pete",
				IsEmailVerified: false,
			},
		}
	}

	return &Mock{data: data}
}

func (m *Mock) FindById(userId int) (*UserDTO, error) {
	for _, user := range m.data {
		if user.Id == userId {
			return user, nil
		}
	}

	return nil, nil
}
