-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL CHECK ( user_id > 0 ),
    extra_information TEXT,
    status smallint NOT NULL,
    verification_code VARCHAR(6)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
