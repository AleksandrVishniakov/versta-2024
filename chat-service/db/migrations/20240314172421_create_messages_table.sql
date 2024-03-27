-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages
(
    id               SERIAL PRIMARY KEY,
    message          TEXT        NOT NULL,
    sender_id        INT         NOT NULL REFERENCES chat_users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    read_by_sender   BOOLEAN     NOT NULL DEFAULT true,
    receiver_id      INT         NOT NULL REFERENCES chat_users (id) ON DELETE CASCADE ON UPDATE CASCADE,
    read_by_receiver BOOLEAN     NOT NULL DEFAULT false,
    created_at       timestamptz NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
