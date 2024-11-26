-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    product_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description TEXT NOT NULL,
    position INT NOT NULL DEFAULT 0,
    img TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
