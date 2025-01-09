-- +goose Up
ALTER TABLE answers
    ADD COLUMN answered_at TIMESTAMP WITHOUT TIME ZONE;

UPDATE answers
set answered_at = now()
where answered_at is null;

-- +goose Down
ALTER TABLE answers
    DROP COLUMN answered_at;
