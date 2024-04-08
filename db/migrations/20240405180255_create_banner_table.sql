-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner
(
    id         bigserial not null primary key,
    name       text      not null,
    feature_id integer   not null references feature on delete cascade ,
    content_id integer not null references content on delete cascade ,
    is_active  boolean   not null default false,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner;
-- +goose StatementEnd
