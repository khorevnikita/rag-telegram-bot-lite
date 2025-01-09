-- +goose Up
CREATE TABLE messages
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users (id),
    content    TEXT,
    response   TEXT,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE messages;
