package usersrepo

import (
	"database/sql"
	"errors"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
)

var (
	ErrUnavailableDatabase = errors.New("database is unavailable")
	ErrUserNotFound        = errors.New("user not found")
)

type UsersRepository interface {
	Create(user *UserEntity) (int, error)

	FindByEmail(email string) (*UserEntity, error)
	FindById(id int) (*UserEntity, error)

	UpdateName(id int, name string) error
	UpdateVerificationCode(id int, code string) error
	MarkEmailAsVerified(id int) error
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) (UsersRepository, error) {
	err := db.Ping()
	if err != nil {
		return nil, wrapErr(ErrUnavailableDatabase)
	}

	return &usersRepository{db: db}, err
}

func (u *usersRepository) Create(user *UserEntity) (id int, err error) {
	defer func() { err = wrapErr(err) }()

	row := u.db.QueryRow(
		`INSERT INTO users (email, name, email_verification_code)
				VALUES ($1, $2, $3)
				RETURNING id`,
		user.Email,
		user.Name,
		user.EmailVerificationCode.String,
	)

	id = 0

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *usersRepository) FindByEmail(email string) (user *UserEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := u.db.QueryRow(
		`SELECT u.* FROM users u
				WHERE u.email = $1`,
		email,
	)

	user = &UserEntity{}

	err = row.Scan(&user.Id, &user.Email, &user.Name, &user.EmailVerificationCode, &user.IsEmailVerified, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *usersRepository) FindById(id int) (user *UserEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := u.db.QueryRow(
		`SELECT u.* FROM users u
				WHERE u.id = $1`,
		id,
	)

	user = &UserEntity{}

	err = row.Scan(&user.Id, &user.Email, &user.Name, &user.EmailVerificationCode, &user.IsEmailVerified, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *usersRepository) UpdateName(id int, name string) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = u.db.Exec(
		`UPDATE users
				SET name=$2
				WHERE id=$1`,
		id,
		name,
	)

	return err
}

func (u *usersRepository) UpdateVerificationCode(id int, code string) (err error) {
	defer func() { err = wrapErr(err) }()

	var verificationCode = &code

	if *verificationCode == "" {
		verificationCode = nil
	}

	_, err = u.db.Exec(
		`UPDATE users
				SET email_verification_code=$2
				WHERE id=$1`,
		id,
		verificationCode,
	)

	return err
}

func (u *usersRepository) MarkEmailAsVerified(id int) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = u.db.Exec(
		`UPDATE users
				SET is_email_verified=true
				WHERE id=$1 AND NOT is_email_verified`,
		id,
	)

	return err
}

func wrapErr(err error) error {
	return e.WrapIfNotNil(err, "usersrepo")
}
