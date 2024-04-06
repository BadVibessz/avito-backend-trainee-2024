-- +goose Up
-- +goose StatementBegin
CREATE TABLE content
(
    id    bigserial not null primary key,
    title text      not null,
    text  text      not null,
    url   text      not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE content;
-- +goose StatementEnd
