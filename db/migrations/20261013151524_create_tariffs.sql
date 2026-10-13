-- +goose Up
CREATE TABLE tariffs (
    id            VARCHAR(36)      PRIMARY KEY,
    title         VARCHAR(255)     NOT NULL,
    num_projects  INTEGER          NOT NULL,
    num_employees INTEGER          NOT NULL,
    price         DOUBLE PRECISION NOT NULL,
    is_removable  BOOLEAN          NOT NULL,
    date_created  BIGINT           NOT NULL,
    date_updated  BIGINT           NOT NULL
);

-- +goose Down
DROP TABLE tariffs;