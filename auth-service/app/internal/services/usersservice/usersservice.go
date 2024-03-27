package usersservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/api/emailapi"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/encryptor"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrMismatchedCodes = errors.New("verification codes do not match")
)

type UsersService interface {
	Register(email string, withEmail bool) (int, error)
	InitAdmin(email string) error

	VerifyEmail(email string, code string) error

	FindById(id int) (*UserResponseDTO, error)
	FindByEmail(email string) (*UserResponseDTO, error)
	FindAdmin() (*UserResponseDTO, error)

	GetVerificationCode(email string) (*VerificationCodeDTO, error)

	UpdateName(id int, name string) error
}

type usersService struct {
	ctx             context.Context
	usersRepository usersrepo.UsersRepository
	emailAPI        emailapi.API
	crypt           *encryptor.Encryptor
}

func NewUsersService(
	ctx context.Context,
	usersRepository usersrepo.UsersRepository,
	emailAPI emailapi.API,
	crypt *encryptor.Encryptor,
) UsersService {
	return &usersService{
		ctx:             ctx,
		usersRepository: usersRepository,
		emailAPI:        emailAPI,
		crypt:           crypt,
	}
}

func (u *usersService) Register(email string, withEmail bool) (int, error) {
	user, err := u.usersRepository.FindByEmail(email)

	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return u.authNewUser(email, withEmail)
	}

	if err != nil {
		return 0, err
	}

	return u.authExistingUser(user, withEmail)
}

func (u *usersService) InitAdmin(email string) error {
	user, err := u.FindAdmin()
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}

	if user != nil {
		return nil
	}

	_, err = u.usersRepository.Create(&usersrepo.UserEntity{
		Email:           email,
		Status:          string(StatusAdmin),
		IsEmailVerified: true,
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *usersService) VerifyEmail(email string, code string) error {
	user, err := u.usersRepository.FindByEmail(email)
	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	var userCode = user.EmailVerificationCode.String
	if userCode != code {
		return ErrMismatchedCodes
	}

	err = u.usersRepository.UpdateVerificationCode(user.Id, "")
	if err != nil {
		return err
	}

	if !user.IsEmailVerified {
		err = u.usersRepository.MarkEmailAsVerified(user.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *usersService) FindById(id int) (*UserResponseDTO, error) {
	entity, err := u.usersRepository.FindById(id)
	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	user, err := mapResponseFromEntity(u.crypt, entity)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersService) FindByEmail(email string) (*UserResponseDTO, error) {
	entity, err := u.usersRepository.FindByEmail(email)
	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	user, err := mapResponseFromEntity(u.crypt, entity)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersService) FindAdmin() (*UserResponseDTO, error) {
	entity, err := u.usersRepository.FindAdmin()
	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	user, err := mapResponseFromEntity(u.crypt, entity)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersService) GetVerificationCode(email string) (*VerificationCodeDTO, error) {
	user, err := u.usersRepository.FindByEmail(email)

	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return nil, ErrUserNotFound
	}

	return &VerificationCodeDTO{
		VerificationCode: user.EmailVerificationCode.String,
	}, nil
}

func (u *usersService) UpdateName(id int, name string) error {
	name, err := encryptString(u.crypt, name)
	if err != nil {
		return err
	}

	err = u.usersRepository.UpdateName(id, name)

	return err
}

func (u *usersService) authNewUser(email string, withEmail bool) (int, error) {
	newUser := &usersrepo.UserEntity{
		Email:                 email,
		Status:                string(StatusUser),
		EmailVerificationCode: sql.NullString{String: newVerificationCode()},
	}

	id, err := u.usersRepository.Create(newUser)
	if err != nil {
		return 0, err
	}

	if !withEmail {
		return id, nil
	}

	user, err := u.usersRepository.FindByEmail(email)
	if err != nil {
		return 0, err
	}

	err = u.sendVerificationCode(email, user.EmailVerificationCode.String)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *usersService) authExistingUser(user *usersrepo.UserEntity, withEmail bool) (int, error) {
	newCode := newVerificationCode()

	err := u.usersRepository.UpdateVerificationCode(user.Id, newCode)
	if err != nil {
		return 0, err
	}

	if !withEmail {
		return user.Id, nil
	}

	err = u.sendVerificationCode(user.Email, newCode)
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (u *usersService) sendVerificationCode(email string, code string) error {
	log.Printf("%s verification code: %s\n", email, code)

	err := u.emailAPI.Write(emailContent(email, code))
	if err != nil {
		return err
	}

	return nil
}

func newVerificationCode() string {
	var digits = "0123456789"

	var code = make([]byte, 6)

	rand.NewSource(time.Now().UnixNano())

	for i := range code {
		var randIndex = rand.Intn(len(digits))

		code[i] = digits[randIndex]
	}

	return string(code)
}

func encryptString(crypt *encryptor.Encryptor, str string) (string, error) {
	encrBytes, err := crypt.Encrypt([]byte(str))

	if err != nil {
		return "", err
	}

	return string(encrBytes), nil
}

func emailContent(email string, code string) *emailapi.EmailDTO {
	return &emailapi.EmailDTO{
		To:      email,
		Subject: fmt.Sprintf("Verification Code [%s]", code),
		Body:    fmt.Sprintf("Thank you for registration! Your email verification code is:\n%s", code),
	}
}
