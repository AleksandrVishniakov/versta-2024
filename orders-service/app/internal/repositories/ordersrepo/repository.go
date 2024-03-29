package ordersrepo

import (
	"database/sql"
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/e"
)

var (
	ErrUnavailableDatabase = errors.New("database is unavailable")
	ErrNoOrders            = errors.New("no orders found")
)

type OrdersRepository interface {
	Create(order *OrderEntity) (int, error)

	FindById(id int, userId int) (*OrderEntity, error)
	FindAll(userId int) ([]*OrderEntity, error)

	UpdateStatus(id int, status byte) error

	Delete(id int, userId int) error
}

type ordersRepository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) (OrdersRepository, error) {
	err := db.Ping()
	if err != nil {
		return nil, wrapErr(ErrUnavailableDatabase)
	}

	return &ordersRepository{db: db}, nil
}

func (o *ordersRepository) Create(order *OrderEntity) (id int, err error) {
	defer func() { err = wrapErr(err) }()

	row := o.db.QueryRow(
		`INSERT INTO orders (user_id, extra_information, status, verification_code)
				VALUES($1, $2, $3, $4)
				RETURNING id`,
		order.UserId,
		order.ExtraInformation.String,
		order.Status,
		order.VerificationCode.String,
	)

	id = 0

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (o *ordersRepository) FindById(id int, userId int) (order *OrderEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := o.db.QueryRow(
		`SELECT * FROM orders o
				WHERE id=$1 AND user_id=$2`,
		id,
		userId,
	)

	order = &OrderEntity{}

	err = row.Scan(&order.Id, &order.UserId, &order.ExtraInformation, &order.Status, &order.VerificationCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoOrders
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *ordersRepository) FindAll(userId int) (orders []*OrderEntity, err error) {
	defer func() { err = wrapErr(err) }()

	rows, err := o.db.Query(
		`SELECT * FROM orders o
				WHERE user_id=$1
				ORDER BY id`,
		userId,
	)

	if err != nil {
		return nil, err
	}

	orders = []*OrderEntity{}

	for rows.Next() {
		var order = &OrderEntity{}

		err = rows.Scan(&order.Id, &order.UserId, &order.ExtraInformation, &order.Status, &order.VerificationCode)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (o *ordersRepository) UpdateStatus(id int, status byte) error {
	_, err := o.db.Exec(
		`UPDATE orders
				SET status=$2
				WHERE id=$1`,
		id,
		status,
	)

	return wrapErr(err)
}

func (o *ordersRepository) Delete(id int, userId int) error {
	_, err := o.db.Exec(
		`DELETE FROM orders
				WHERE id=$1 AND user_id=$2`,
		id,
		userId,
	)

	return wrapErr(err)
}

func wrapErr(err error) error {
	return e.WrapIfNotNil(err, "ordersrepo")
}
