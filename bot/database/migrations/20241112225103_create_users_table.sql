-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users
(
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id       BIGINT UNIQUE,
    telegram_username TEXT,
    first_name        TEXT,
    last_name         TEXT,
    avatar            TEXT,
    connection_date   TIMESTAMP,
    message_count     INT              DEFAULT 0,
    last_message_date TIMESTAMP,
    conversation_id   INT,
    is_admin          BOOLEAN,
    created_at        TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
