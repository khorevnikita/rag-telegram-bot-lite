-- +goose Up
-- +goose StatementBegin
CREATE TABLE message_files
(
    id         UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    message_id UUID    NOT NULL,
    user_id    UUID    NOT NULL,
    file_path  TEXT    NOT NULL,
    file_name  VARCHAR NOT NULL,
    size       INT     NOT NULL,
    extension  VARCHAR NOT NULL,
    file_type  VARCHAR NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (message_id) REFERENCES messages (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS message_files;
-- +goose StatementEnd
