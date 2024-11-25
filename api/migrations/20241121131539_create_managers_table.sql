-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS managers (
    manager_id BIGINT NOT NULL PRIMARY KEY,
    phone VARCHAR(18) NOT NULL,
    email VARCHAR(255) NOT NULL,
    state INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS managers;
-- +goose StatementEnd
