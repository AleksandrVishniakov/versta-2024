package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/api/authapi"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/internal/api/emailapi"
	"github.com/AleksnadrVishniakov/versta-2024/orders-service/app/pkg/e"
	"log/slog"
	"math/rand"
	"time"
)

var (
	ErrWithAuthorization = errors.New("authorization error")
	ErrMismatchedCodes   = errors.New("mismatched codes")
)

type Service interface {
	Create(sessionKey string, email string, extraInformation string) (int, string, error)

	FindAll(sessionKey string) ([]*OrderDTO, string, error)
	FindById(sessionKey string, orderId int) (*OrderDTO, string, error)

	Verify(email string, orderId int, verificationCode string) (string, error)
	Complete(orderId int) error

	Delete(sessionKey string, orderId int) (string, error)
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
	id          int
	email       string
	name        string
	sessionKey  string
	withSession bool
}

func (o *ordersService) Create(sessionKey string, email string, extraInformation string) (int, string, error) {
	user, err := o.authOrCreateUser(sessionKey, email)
	if err != nil {
		return 0, "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	order := &OrderDTO{
		UserId:           user.id,
		ExtraInformation: extraInformation,
		Status:           StatusCreated,
	}
	orderVerificationCode := verificationCode()

	orderId, err := o.storage.Create(order, orderVerificationCode)
	if err != nil {
		return 0, "", err
	}

	o.sendVerificationMessage(
		user.email,
		user.name,
		orderId,
		extraInformation,
		orderVerificationCode,
	)

	return orderId, user.sessionKey, nil
}

func (o *ordersService) FindAll(sessionKey string) ([]*OrderDTO, string, error) {
	user, err := o.authUser(sessionKey)
	if err != nil {
		return nil, "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	orders, err := o.storage.FindAll(user.id)
	if err != nil {
		return nil, "", err
	}

	return orders, user.sessionKey, nil
}

func (o *ordersService) FindById(sessionKey string, orderId int) (*OrderDTO, string, error) {
	user, err := o.authUser(sessionKey)
	if err != nil {
		return nil, "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	order, err := o.storage.FindById(orderId, user.id)
	if err != nil {
		return nil, "", err
	}

	return order, user.sessionKey, nil
}

func (o *ordersService) Verify(email string, orderId int, verificationCode string) (string, error) {
	userDTO, err := o.authAPI.FindByEmail(email)
	if err != nil {
		return "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	orderVerificationCode, err := o.storage.GetVerificationCode(orderId, userDTO.Id)
	if err != nil {
		return "", err
	}

	if orderVerificationCode != verificationCode {
		return "", ErrMismatchedCodes
	}

	err = o.storage.MarkAsVerified(orderId)
	if err != nil {
		return "", err
	}

	sessionKey, err := o.verifyUser(email)
	if err != nil {
		return "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	return sessionKey, nil
}

func (o *ordersService) Complete(orderId int) error {
	err := o.storage.MarkAsCompleted(orderId)
	if err != nil {
		return err
	}

	return nil
}

func (o *ordersService) Delete(sessionKey string, orderId int) (string, error) {
	user, err := o.authUser(sessionKey)
	if err != nil {
		return "", e.WrapErrWithErr(err, ErrWithAuthorization)
	}

	err = o.storage.Delete(orderId, user.id)
	if err != nil {
		return "", err
	}

	return user.sessionKey, nil
}

func (o *ordersService) authOrCreateUser(sessionKey string, email string) (*userData, error) {
	var err error
	var user = &userData{
		email:       email,
		sessionKey:  sessionKey,
		withSession: sessionKey != "",
	}

	if !user.withSession {
		id, err := o.authAPI.Create(user.email, false)
		if err != nil {
			return nil, err
		}

		userDTO, err := o.authAPI.FindByEmail(email)
		if err != nil {
			return nil, err
		}

		user.id = id
		user.name = userDTO.Name
	} else {
		user, err = o.authUser(sessionKey)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (o *ordersService) authUser(sessionKey string) (*userData, error) {
	var user = &userData{
		withSession: true,
		sessionKey:  sessionKey,
	}

	userDTO, newSessionKey, err := o.authAPI.FindBySessionKey(user.sessionKey)
	if err != nil {
		return nil, err
	}

	user.id = userDTO.Id
	user.email = userDTO.Email
	user.name = userDTO.Name
	user.sessionKey = newSessionKey

	return user, nil
}

func (o *ordersService) verifyUser(email string) (string, error) {
	vCode, err := o.authAPI.GetVerificationCode(email)
	if err != nil {
		return "", err
	}

	sessionKey, err := o.authAPI.VerifyEmail(email, vCode.VerificationCode)
	if err != nil {
		return "", err
	}

	return sessionKey, nil
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
