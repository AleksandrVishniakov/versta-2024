package chattersrepo

import (
	"database/sql"
	"errors"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
)

var (
	ErrChatterNotFound = errors.New("chattersrepo: chatter not found")
)

type ChattersRepository interface {
	Create(chatter *ChatterEntity) (int, error)

	FindBySession(session string) (*ChatterEntity, error)
	FindByChatterId(chatterId int) (*ChatterEntity, error)
	FindByUserId(userId int) (*ChatterEntity, error)
	FindSendersByChatterId(chatterId int) ([]*ChatterEntityWithUnread, error)

	ChangeSessionToId(session string, userId int) error
}

type chatterRepo struct {
	db *sql.DB
}

func NewChattersRepository(db *sql.DB) ChattersRepository {
	return &chatterRepo{db: db}
}

func (c *chatterRepo) Create(chatter *ChatterEntity) (id int, err error) {
	defer func() { err = wrapErr(err) }()

	row := c.db.QueryRow(
		`INSERT INTO chat_users (user_id, temp_session)
				VALUES ($1, $2) RETURNING id`,
		chatter.UserId.Int32,
		chatter.TempSession.String,
	)

	id = 0

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (c *chatterRepo) FindBySession(session string) (chatter *ChatterEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := c.db.QueryRow(
		`SELECT * FROM chat_users
				WHERE temp_session=$1`,
		session,
	)

	chatter = &ChatterEntity{}

	err = row.Scan(&chatter.Id, &chatter.UserId, &chatter.TempSession)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return chatter, err
}

func (c *chatterRepo) FindByChatterId(chatterId int) (chatter *ChatterEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := c.db.QueryRow(
		`SELECT * FROM chat_users
				WHERE id=$1`,
		chatterId,
	)

	chatter = &ChatterEntity{}

	err = row.Scan(&chatter.Id, &chatter.UserId, &chatter.TempSession)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return chatter, err
}

func (c *chatterRepo) FindByUserId(userId int) (chatter *ChatterEntity, err error) {
	defer func() { err = wrapErr(err) }()

	row := c.db.QueryRow(
		`SELECT * FROM chat_users
				WHERE user_id=$1`,
		userId,
	)

	chatter = &ChatterEntity{}

	err = row.Scan(&chatter.Id, &chatter.UserId, &chatter.TempSession)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrChatterNotFound
	}
	if err != nil {
		return nil, err
	}

	return chatter, err
}

func (c *chatterRepo) FindSendersByChatterId(chatterId int) (chatters []*ChatterEntityWithUnread, err error) {
	defer func() { err = wrapErr(err) }()

	rows, err := c.db.Query(
		`SELECT DISTINCT c.*, (
    				SELECT COUNT(*) as count
    				FROM messages
    				WHERE receiver_id=$1 AND sender_id=c.id AND NOT read_by_receiver
				) as unread_messages_count  FROM chat_users c
				LEFT JOIN messages m ON m.sender_id = c.id
				WHERE m.receiver_id=$1`,
		chatterId,
	)

	if err != nil {
		return nil, err
	}

	chatters = []*ChatterEntityWithUnread{}

	for rows.Next() {
		chatter := &ChatterEntityWithUnread{}

		err := rows.Scan(&chatter.Id, &chatter.UserId, &chatter.TempSession, &chatter.UnreadMessagesCount)
		if err != nil {
			return nil, err
		}

		chatters = append(chatters, chatter)
	}

	return chatters, nil
}

func (c *chatterRepo) ChangeSessionToId(session string, userId int) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = c.db.Exec(
		`UPDATE chat_users
				SET user_id=$1,
				    temp_session=null
				WHERE temp_session=$2`,
		userId,
		session,
	)

	return err
}

func wrapErr(err error) error {
	return e.WrapIfNotNil(err, "chat_users_repo")
}
