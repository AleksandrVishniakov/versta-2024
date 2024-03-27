-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_users(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    temp_session VARCHAR(16)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chat_users;
-- +goose StatementEnd
