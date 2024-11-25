-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    product_id BIGINT NOT NULL PRIMARY KEY,
    description TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
