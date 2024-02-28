package sessionsrepo

import (
	"database/sql"
	"errors"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
)

var (
	ErrUnavailableDatabase = errors.New("database is unavailable")
)

type SessionsRepository interface {
	Delete(id int) error
	Create(session *SessionEntity) error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionsRepository(db *sql.DB) (SessionsRepository, error) {
	err := db.Ping()
	if err != nil {
		return nil, wrapErr(ErrUnavailableDatabase)
	}

	return &sessionRepository{db: db}, nil
}

func (s *sessionRepository) Delete(id int) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = s.db.Exec(
		`DELETE FROM sessions
				WHERE id=$1`,
		id,
	)

	return err
}

func (s *sessionRepository) Create(session *SessionEntity) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = s.db.Exec(
		`INSERT INTO sessions (user_id, session_key, created_at, expires_at)
				VALUES ($1, $2, $3, $4)`,
		session.UserId,
		session.SessionKey,
		session.CreatedAt,
		session.ExpiresAt,
	)

	return err
}

func wrapErr(err error) error {
	return e.WrapIfNotNil(err, "sessionsrepo")
}
