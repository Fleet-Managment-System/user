-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE Table "users" (
    id serial primary key,
    firstname VARCHAR(32) not NULL,
    lastname VARCHAR(32) not NULL,
    email VARCHAR(128) not NULL UNIQUE,
    passwordHash VARCHAR(32) not NULL,
    createdAt timestamp not NULL DEFAULT now(),
    updatedAt timestamp
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "user";
-- +goose StatementEnd