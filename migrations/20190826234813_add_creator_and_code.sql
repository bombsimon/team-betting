-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE competition
    ADD COLUMN code varchar(10) AFTER description,
    ADD COLUMN created_by_id INT,
    ADD CONSTRAINT competition_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES better(id) ON DELETE CASCADE;

ALTER TABLE competitor
    ADD COLUMN created_by_id INT,
    ADD CONSTRAINT competitor_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES better(id) ON DELETE CASCADE;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE competition
    DROP FOREIGN KEY competition_created_by_id_fk,
    DROP COLUMN created_by_id,
    DROP COLUMN code;

ALTER TABLE competitor
    DROP FOREIGN KEY competition_created_by_id_fk,
    DROP COLUMN created_by_id;
