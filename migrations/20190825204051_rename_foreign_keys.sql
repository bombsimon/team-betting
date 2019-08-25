-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE `competition_competitor`
    CHANGE `id_competition` `competition_id` INT NOT NULL,
    CHANGE `id_competitor` `competitor_id` INT NOT NULL;

ALTER TABLE `bet`
    CHANGE `id_better` `better_id` INT NOT NULL,
    CHANGE `id_competition_competitor` `competition_competitor_id` INT NOT NULL;

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

ALTER TABLE `competition_competitor`
    CHANGE `competition_id` `id_competition` INT NOT NULL,
    CHANGE `competitor_id` `id_competitor` INT NOT NULL;

ALTER TABLE `bet`
    CHANGE `better_id` `id_better` INT NOT NULL,
    CHANGE `competition_competitor_id` `id_competition_competitor` INT NOT NULL;
