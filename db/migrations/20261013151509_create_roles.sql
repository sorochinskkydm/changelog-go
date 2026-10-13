-- +goose Up
CREATE TABLE roles (
    id    VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE roles;