package messagesrepo

import (
	"database/sql"
	"github.com/AleksandrVishniakov/versta-2024/chat-service/app/pkg/e"
)

type MessagesRepository interface {
	Create(message *MessageEntity) (int, error)

	FindByChatterId(chatterId int) ([]*MessageEntity, error)
	FindBySenderAndReceiver(senderId, receiverId int) ([]*MessageEntity, error)

	GetUnreadCount(forId, withId int) (int, error)

	ReadAll(forId, withId int) error
}

type msgRepo struct {
	db *sql.DB
}

func NewMessagesRepository(db *sql.DB) MessagesRepository {
	return &msgRepo{db: db}
}

func (m *msgRepo) Create(message *MessageEntity) (id int, err error) {
	defer func() { err = wrapErr(err) }()

	row := m.db.QueryRow(
		`INSERT INTO messages (message, sender_id, receiver_id, read_by_receiver)
				VALUES ($1, $2, $3, $4) RETURNING id`,
		message.Message,
		message.SenderId,
		message.ReceiverId,
		message.ReadByReceiver,
	)

	id = 0

	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *msgRepo) FindByChatterId(chatterId int) (messages []*MessageEntity, err error) {
	defer func() { err = wrapErr(err) }()

	rows, err := m.db.Query(
		`SELECT * FROM messages
				WHERE sender_id=$1 OR receiver_id=$1
				ORDER BY created_at`,
		chatterId,
	)

	if err != nil {
		return nil, err
	}

	messages = []*MessageEntity{}

	messages, err = m.parseMessages(rows, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *msgRepo) FindBySenderAndReceiver(senderId, receiverId int) (messages []*MessageEntity, err error) {
	defer func() { err = wrapErr(err) }()

	rows, err := m.db.Query(
		`SELECT * FROM messages
				WHERE 
				    (sender_id=$1 AND receiver_id=$2) OR (receiver_id=$1 AND sender_id=$2)
				ORDER BY created_at`,
		senderId,
		receiverId,
	)

	if err != nil {
		return nil, err
	}

	messages = []*MessageEntity{}

	messages, err = m.parseMessages(rows, messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *msgRepo) GetUnreadCount(forId, withId int) (count int, err error) {
	defer func() { err = wrapErr(err) }()

	row := m.db.QueryRow(
		`SELECT COUNT(*) as count
				FROM messages
				WHERE receiver_id=$1 AND sender_id=$2 AND NOT read_by_receiver`,
		forId,
		withId,
	)

	count = 0

	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *msgRepo) ReadAll(forId, withId int) (err error) {
	defer func() { err = wrapErr(err) }()

	_, err = m.db.Exec(
		`UPDATE messages
				SET read_by_receiver=true
				WHERE receiver_id=$1 AND sender_id=$2 AND NOT read_by_receiver`,
		forId,
		withId,
	)

	return err
}

func (m *msgRepo) parseMessages(rows *sql.Rows, messages []*MessageEntity) ([]*MessageEntity, error) {
	for rows.Next() {
		msg := &MessageEntity{}

		err := rows.Scan(
			&msg.Id,
			&msg.Message,
			&msg.SenderId,
			&msg.ReadBySender,
			&msg.ReceiverId,
			&msg.ReadByReceiver,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}
	return messages, nil
}

func wrapErr(err error) error {
	return e.WrapIfNotNil(err, "messages_repo")
}
