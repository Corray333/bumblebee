-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    date BIGINT NOT NULL,
    customer_phone VARCHAR(18) NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    customer_address VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
