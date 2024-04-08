-- +goose Up
-- +goose StatementBegin
CREATE TABLE banner_tag
(
    id        bigserial not null primary key,
    banner_id integer   not null references banner on delete cascade,
    tag_id    integer   not null references tag
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE banner_tag
-- +goose StatementEnd
