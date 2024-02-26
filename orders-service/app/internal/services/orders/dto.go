package orders

import "fmt"

type OrderDTO struct {
	Id               int    `json:"id"`
	UserId           int    `json:"userId"`
	ExtraInformation string `json:"extraInformation"`
	Status           OrderStatus
}

func (o *OrderDTO) Valid() (bool, error) {
	if o.Id < 1 {
		return false, fmt.Errorf("id cannot be less than one, got id=%d", o.Id)
	}

	if o.UserId < 1 {
		return false, fmt.Errorf("user_id cannot be less than one, got user_id=%d", o.UserId)
	}

	if o.Status < 0 {
		return false, fmt.Errorf("unknown status %d", o.Status)
	}

	return true, nil
}
