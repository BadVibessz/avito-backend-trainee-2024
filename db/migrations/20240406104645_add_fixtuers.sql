-- +goose Up
-- +goose StatementBegin
INSERT INTO tag (name) VALUES ('tag_test_name');
INSERT INTO tag (name) VALUES ('tag_test_name2');


INSERT INTO feature (name) VALUES ('feature_test_name');
INSERT INTO feature (name) VALUES ('feature_test_name2');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
