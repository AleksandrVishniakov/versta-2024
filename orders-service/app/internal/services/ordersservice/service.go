package ordersservice

import (
	"errors"

	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/repositories/ordersrepo"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/orders"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/scrambler"
)

var ErrNoOrders = errors.New("no orders found")

type OrdersService interface {
	Create(order *orders.OrderDTO) (int, error)

	FindById(id int, userId int) (*orders.OrderDTO, error)
	FindAll(userId int) ([]*orders.OrderDTO, error)

	MarkAsVerified(id int) error
	MarkAsCompleted(id int) error

	Delete(id int, userId int) error
}

type ordersService struct {
	repository ordersrepo.OrdersRepository
	scrambler  scrambler.Scrambler
}

func NewOrdersService(
	repo ordersrepo.OrdersRepository,
	scrambler scrambler.Scrambler,
) OrdersService {
	return &ordersService{
		repository: repo,
		scrambler:  scrambler,
	}
}

func (o *ordersService) Create(order *orders.OrderDTO) (int, error) {
	order.Status = orders.StatusCreated

	entity, err := orders.MapEntityFromDTO(order, o.scrambler)
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

func (o *ordersService) FindById(id int, userId int) (*orders.OrderDTO, error) {
	entity, err := o.repository.FindById(id, userId)
	if errors.Is(err, ordersrepo.ErrNoOrders) {
		return nil, ErrNoOrders
	}

	if err != nil {
		return nil, err
	}

	order, err := orders.MapDTOFromEntity(entity, o.scrambler)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *ordersService) FindAll(userId int) ([]*orders.OrderDTO, error) {
	entities, err := o.repository.FindAll(userId)
	if err != nil {
		return nil, err
	}

	var userOrders []*orders.OrderDTO

	for _, e := range entities {
		order, err := orders.MapDTOFromEntity(e, o.scrambler)
		if err != nil {
			return nil, err
		}

		userOrders = append(userOrders, order)
	}

	if len(userOrders) == 0 {
		return []*orders.OrderDTO{}, nil
	}

	return userOrders, nil
}

func (o *ordersService) MarkAsVerified(id int) error {
	err := o.repository.UpdateStatus(id, byte(orders.StatusVerified))
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) MarkAsCompleted(id int) error {
	err := o.repository.UpdateStatus(id, byte(orders.StatusCompleted))
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) Delete(id int, userId int) error {
	err := o.repository.Delete(id, userId)
	if err != nil {
		return err
	}

	return nil
}
