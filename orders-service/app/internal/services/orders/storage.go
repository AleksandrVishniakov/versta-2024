package orders

import (
	"errors"

	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/internal/repositories/ordersrepo"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/scrambler"
)

var ErrNoOrders = errors.New("no orders found")

type Storage interface {
	Create(order *OrderDTO, verificationCode string) (int, error)

	FindById(id int, userId int) (*OrderDTO, error)
	FindAll(userId int) ([]*OrderDTO, error)

	GetVerificationCode(id int, userId int) (string, error)

	MarkAsVerified(id int) error
	MarkAsCompleted(id int) error

	Delete(id int, userId int) error
}

type storage struct {
	repository ordersrepo.OrdersRepository
	scrambler  scrambler.Scrambler
}

func NewOrdersStorage(
	repo ordersrepo.OrdersRepository,
	scrambler scrambler.Scrambler,
) Storage {
	return &storage{
		repository: repo,
		scrambler:  scrambler,
	}
}

func (o *storage) Create(order *OrderDTO, verificationCode string) (int, error) {
	order.Status = StatusCreated

	entity, err := MapEntityFromDTO(order, verificationCode, o.scrambler)
	if err != nil {
		return 0, err
	}

	id, err := o.repository.Create(
		entity,
	)

	if err != nil {
		return 0, err
	}

	return id, err
}

func (o *storage) FindById(id int, userId int) (*OrderDTO, error) {
	entity, err := o.repository.FindById(id, userId)
	if errors.Is(err, ordersrepo.ErrNoOrders) {
		return nil, ErrNoOrders
	}

	if err != nil {
		return nil, err
	}

	order, err := MapDTOFromEntity(entity, o.scrambler)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *storage) FindAll(userId int) ([]*OrderDTO, error) {
	entities, err := o.repository.FindAll(userId)
	if err != nil {
		return nil, err
	}

	var userOrders []*OrderDTO

	for _, e := range entities {
		order, err := MapDTOFromEntity(e, o.scrambler)
		if err != nil {
			return nil, err
		}

		userOrders = append(userOrders, order)
	}

	if len(userOrders) == 0 {
		return []*OrderDTO{}, nil
	}

	return userOrders, nil
}

func (o *storage) GetVerificationCode(id int, userId int) (string, error) {
	entity, err := o.repository.FindById(id, userId)
	if errors.Is(err, ordersrepo.ErrNoOrders) {
		return "", ErrNoOrders
	}

	if err != nil {
		return "", err
	}

	return entity.VerificationCode.String, nil
}

func (o *storage) MarkAsVerified(id int) error {
	err := o.repository.UpdateStatus(id, byte(StatusVerified))
	if err != nil {
		return err
	}

	return nil
}

func (o *storage) MarkAsCompleted(id int) error {
	err := o.repository.UpdateStatus(id, byte(StatusCompleted))
	if err != nil {
		return err
	}

	return nil
}

func (o *storage) Delete(id int, userId int) error {
	err := o.repository.Delete(id, userId)
	if err != nil {
		return err
	}

	return nil
}
