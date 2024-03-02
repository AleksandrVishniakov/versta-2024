package sessionsrepo

import (
	"database/sql"
	"errors"

	"github.com/AleksandrVishniakov/versta-2024/auth-service/app/pkg/e"
)

var (
	ErrUnavailableDatabase = errors.New("database is unavailable")
	ErrSessionNotFound     = errors.New("session not found")
)

type SessionsRepository interface {
	Create(session *SessionEntity) error

	FindByKey(sessionKey string) (*SessionEntity, error)

	Delete(id int) error
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

func (s *sessionRepository) FindByKey(sessionKey string) (session *SessionEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := s.db.QueryRow(
		`SELECT s.* FROM sessions s
				WHERE s.session_key = $1`,
		sessionKey,
	)

	session = &SessionEntity{}

	err = row.Scan(&session.Id, &session.UserId, &session.SessionKey, &session.CreatedAt, &session.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return session, nil
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
