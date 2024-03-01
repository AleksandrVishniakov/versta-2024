package usersservice

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/sessionsrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/internal/repositories/usersrepo"
	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/encryptor"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrMismatchedCodes = errors.New("verification codes do not match")
)

type UsersService interface {
	Register(email string) (int, error)

	VerifyEmail(email string, code string) error

	FindBySessionKey(sessionKey string) (*UserResponseDTO, error)
	GetVerificationCode(email string) (*VerificationCodeDTO, error)

	UpdateName(id int, name string) error
}

type usersService struct {
	ctx                context.Context
	usersRepository    usersrepo.UsersRepository
	sessionsRepository sessionsrepo.SessionsRepository
	crypt              *encryptor.Encryptor
}

func NewUsersService(
	ctx context.Context,
	usersRepository usersrepo.UsersRepository,
	sessionsRepository sessionsrepo.SessionsRepository,
	crypt *encryptor.Encryptor,
) UsersService {
	return &usersService{
		ctx:                ctx,
		usersRepository:    usersRepository,
		sessionsRepository: sessionsRepository,
		crypt:              crypt,
	}
}

func (u *usersService) Register(email string) (int, error) {
	email, err := encryptString(u.crypt, email)
	if err != nil {
		return 0, err
	}

	user, err := u.usersRepository.FindByEmail(email)

	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return u.authNewUser(email)
	}

	if err != nil {
		return 0, err
	}

	return u.authExistingUser(user)
}

func (u *usersService) VerifyEmail(email string, code string) error {
	email, err := encryptString(u.crypt, email)
	if err != nil {
		return err
	}

	user, err := u.usersRepository.FindByEmail(email)
	if errors.Is(err, usersrepo.ErrUserNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return err
	}

	if user.IsEmailVerified {
		return nil
	}

	var userCode = user.EmailVerificationCode.String
	if userCode != code {
		return ErrMismatchedCodes
	}

	err = u.usersRepository.UpdateVerificationCode(user.Id, "")
	if err != nil {
		return err
	}

	err = u.usersRepository.MarkEmailAsVerified(user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u *usersService) FindBySessionKey(sessionKey string) (*UserResponseDTO, error) {
	entity, err := u.usersRepository.FindBySessionKey(sessionKey)
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
	email, err := encryptString(u.crypt, email)
	if err != nil {
		return nil, err
	}

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

func (u *usersService) authNewUser(encryptedEmail string) (int, error) {
	newUser := &usersrepo.UserEntity{
		Email:                 encryptedEmail,
		EmailVerificationCode: sql.NullString{String: newVerificationCode()},
	}

	id, err := u.usersRepository.Create(newUser)
	if err != nil {
		return 0, nil
	}

	user, err := u.usersRepository.FindByEmail(encryptedEmail)
	if err != nil {
		return 0, nil
	}

	email, err := decryptString(u.crypt, user.Email)
	if err != nil {
		return 0, nil
	}

	err = u.sendVerificationCode(email, user.EmailVerificationCode.String)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

func (u *usersService) authExistingUser(user *usersrepo.UserEntity) (int, error) {
	newCode := newVerificationCode()

	err := u.usersRepository.UpdateVerificationCode(user.Id, newCode)
	if err != nil {
		return 0, err
	}

	email, err := decryptString(u.crypt, user.Email)
	if err != nil {
		return 0, nil
	}

	err = u.sendVerificationCode(email, newCode)
	if err != nil {
		return 0, nil
	}

	return user.Id, nil
}

func (u *usersService) sendVerificationCode(email string, code string) error {
	// TODO: implement sending email logic
	log.Printf("%s verification code: %s\n", email, code)
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

func decryptString(crypt *encryptor.Encryptor, str string) (string, error) {
	decrBytes, err := crypt.Decrypt([]byte(str))

	if err != nil {
		return "", err
	}

	return string(decrBytes), nil
}
