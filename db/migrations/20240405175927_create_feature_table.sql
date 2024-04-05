-- +goose Up
-- +goose StatementBegin
CREATE TABLE feature
(
    id bigserial not null primary key,
    name text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE feature;
-- +goose StatementEnd
