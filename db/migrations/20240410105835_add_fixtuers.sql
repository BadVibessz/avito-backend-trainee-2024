-- +goose Up
-- +goose StatementBegin
INSERT INTO tag (name) VALUES ('tag_test_name');
INSERT INTO tag (name) VALUES ('tag_test_name2');


INSERT INTO feature (name) VALUES ('feature_test_name');
INSERT INTO feature (name) VALUES ('feature_test_name2');

INSERT INTO users (username, is_admin, hashed_password) VALUES ('admin', true, '$2a$10$2tfCkk9Gtr4Aege5b.u9/e8BtqgrpGEgh/2WkLxHoNO3aeQhf1Z/G');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE tag CASCADE;
TRUNCATE TABLE feature CASCADE;
TRUNCATE TABLE users;
-- +goose StatementEnd
