package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/internal/api/authapi"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/internal/api/emailapi"
	"github.com/AleksandrVishniakov/versta-2024/orders-service/app/pkg/e"
	"log/slog"
	"math/rand"
	"time"
)

var (
	ErrWithAuthorization = errors.New("authorization error")
	ErrMismatchedCodes   = errors.New("mismatched codes")
)

type Service interface {
	Create(email string, extraInformation string) (int, error)

	FindAll(email string) ([]*OrderDTO, error)
	FindById(email string, orderId int) (*OrderDTO, error)

	Verify(email string, orderId int, verificationCode string) (accessToken string, refreshToken string, err error)
	Complete(orderId int) error

	Delete(email string, orderId int) error
}

type ordersService struct {
	ctx      context.Context
	storage  Storage
	authAPI  authapi.API
	emailAPI emailapi.API
}

func NewOrdersService(
	ctx context.Context,
	storage Storage,
	authAPI authapi.API,
	emailAPI emailapi.API,
) Service {
	return &ordersService{
		ctx:      ctx,
		storage:  storage,
		authAPI:  authAPI,
		emailAPI: emailAPI,
	}
}

type userData struct {
	id    int
	email string
	name  string
}

func (o *ordersService) Create(email string, extraInformation string) (int, error) {
	user, err := o.authOrCreateUser(email)
	if err != nil {
		return 0, e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	order := &OrderDTO{
		UserId:           user.id,
		ExtraInformation: extraInformation,
		Status:           StatusCreated,
	}
	orderVerificationCode := verificationCode()

	orderId, err := o.storage.Create(order, orderVerificationCode)
	if err != nil {
		return 0, err
	}

	o.sendVerificationMessage(
		user.email,
		user.name,
		orderId,
		extraInformation,
		orderVerificationCode,
	)

	return orderId, nil
}

func (o *ordersService) FindAll(email string) ([]*OrderDTO, error) {
	user, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	orders, err := o.storage.FindAll(user.Id)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *ordersService) FindById(email string, orderId int) (*OrderDTO, error) {
	user, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return nil, e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	order, err := o.storage.FindById(orderId, user.Id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *ordersService) Verify(email string, orderId int, verificationCode string) (string, string, error) {
	userDTO, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return "", "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	orderVerificationCode, err := o.storage.GetVerificationCode(orderId, userDTO.Id)
	if err != nil {
		return "", "", err
	}

	if orderVerificationCode != verificationCode {
		return "", "", ErrMismatchedCodes
	}

	err = o.storage.MarkAsVerified(orderId)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := o.verifyUser(email)
	if err != nil {
		return "", "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	return accessToken, refreshToken, nil
}

func (o *ordersService) Complete(orderId int) error {
	err := o.storage.MarkAsCompleted(orderId)
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) Delete(email string, orderId int) error {
	user, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	err = o.storage.Delete(orderId, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) authOrCreateUser(email string) (*userData, error) {
	var err error
	var user = &userData{
		email: email,
	}

	id, err := o.authAPI.Register(user.email, false)
	if err != nil {
		return nil, err
	}

	userDTO, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	user.id = id
	user.name = userDTO.Name

	return user, nil
}

func (o *ordersService) verifyUser(email string) (string, string, error) {
	_, err := o.authAPI.Register(email, false)
	if err != nil {
		return "", "", err
	}

	vCode, err := o.authAPI.GetVerificationCode(email)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := o.authAPI.VerifyEmail(email, vCode.VerificationCode)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (o *ordersService) sendVerificationMessage(
	email string,
	name string,
	orderId int,
	extraOrderInformation string,
	orderVerificationCode string,
) {

	var greeting string
	if name == "" {
		greeting = "Hello!"
	} else {
		greeting = fmt.Sprintf("Hello, %s!", name)
	}

	emailContent := &emailapi.EmailDTO{
		To:      email,
		Subject: "Order verification",
		Body: fmt.Sprintf(
			"\n\n%s You've just made an order #%d with following information:\n\n"+
				"%s\n\n"+
				"If you didn't order this, delete this this email.\n"+
				"Enter this code to verify your order. Don't share your code\n\n"+
				"%s",
			greeting,
			orderId,
			extraOrderInformation,
			orderVerificationCode,
		),
	}

	go func() {
		err := o.emailAPI.Write(emailContent)
		if err != nil {
			slog.Error("email send error",
				slog.String("error", err.Error()),
				slog.String("email", email),
				slog.Int("orderId", orderId),
				slog.String("verificationCode", orderVerificationCode),
			)
		}
	}()
}

func verificationCode() string {
	var digits = "0123456789"

	var code = make([]byte, 6)

	rand.NewSource(time.Now().UnixNano())

	for i := range code {
		var randIndex = rand.Intn(len(digits))

		code[i] = digits[randIndex]
	}

	return string(code)
}
