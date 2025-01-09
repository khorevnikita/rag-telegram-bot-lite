-- +goose Up
ALTER TABLE messages
    ADD COLUMN liked               BOOLEAN,
    ADD COLUMN telegram_message_id INTEGER,
    ADD COLUMN ai_message_id       INTEGER;

-- +goose Down
ALTER TABLE messages
    DROP COLUMN liked,
    DROP COLUMN telegram_message_id,
    DROP COLUMN ai_message_id;
