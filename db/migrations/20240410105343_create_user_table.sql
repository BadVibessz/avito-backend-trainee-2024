-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id              bigserial    not null primary key,
    username        varchar(256) not null,
    is_admin        boolean      not null,
    hashed_password text         not null,
    created_at      timestamp    not null default now(),
    updated_at      timestamp    not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
