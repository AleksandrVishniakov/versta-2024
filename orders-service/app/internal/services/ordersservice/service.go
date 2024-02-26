package ordersservice

import (
	"errors"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/repositories/ordersrepo"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/services/orders"
)

var ErrNoOrders = errors.New("no orders found")

type OrdersService interface {
	Create(order *orders.OrderDTO) (int, error)

	FindById(id int, userId int) (*orders.OrderDTO, error)
	FindAll(userId int) ([]*orders.OrderDTO, error)

	MarkAsVerified(id int, userId int) error
	MarkAsCompleted(id int, userId int) error
}

type ordersService struct {
	repository ordersrepo.OrdersRepository
}

func NewOrdersService(repo ordersrepo.OrdersRepository) OrdersService {
	return &ordersService{repository: repo}
}

func (o *ordersService) Create(order *orders.OrderDTO) (int, error) {
	id, err := o.repository.Create(
		orders.MapEntityFromDTO(order),
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

	return orders.MapDTOFromEntity(entity), nil
}

func (o *ordersService) FindAll(userId int) ([]*orders.OrderDTO, error) {
	entities, err := o.repository.FindAll(userId)
	if err != nil {
		return nil, err
	}

	var userOrders []*orders.OrderDTO

	for _, e := range entities {
		userOrders = append(userOrders, orders.MapDTOFromEntity(e))
	}

	if len(userOrders) == 0 {
		return []*orders.OrderDTO{}, nil
	}

	return userOrders, nil
}

func (o *ordersService) MarkAsVerified(id int, userId int) error {
	err := o.repository.UpdateStatus(id, userId, byte(orders.StatusVerified))
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) MarkAsCompleted(id int, userId int) error {
	err := o.repository.UpdateStatus(id, userId, byte(orders.StatusCompleted))
	if err != nil {
		return err
	}

	return nil
}
