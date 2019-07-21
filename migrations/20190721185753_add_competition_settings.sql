-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE competition
    ADD COLUMN min_score INT NOT NULL DEFAULT 0,
    ADD COLUMN max_score INT NOT NULL DEFAULT 10;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE competition
    DROP COLUMN min_score,
    DROP COLUMN max_score;
