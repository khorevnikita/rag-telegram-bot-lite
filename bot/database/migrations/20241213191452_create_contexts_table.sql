-- +goose Up
-- +goose StatementBegin
CREATE TABLE system_contexts
(
    id         UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    key        VARCHAR NOT NULL,
    text       TEXT    NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at TIMESTAMP WITHOUT TIME ZONE
);

ALTER TABLE users
    ADD COLUMN system_context_id UUID,
    ADD CONSTRAINT fk_users_system_contexts FOREIGN KEY (system_context_id) REFERENCES system_contexts (id)
;

ALTER TABLE messages
    ADD COLUMN system_context_id UUID,
    ADD CONSTRAINT fk_messages_system_contexts FOREIGN KEY (system_context_id) REFERENCES system_contexts (id)
;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages
    DROP CONSTRAINT IF EXISTS fk_messages_system_contexts;
ALTER TABLE messages
    DROP COLUMN IF EXISTS system_context_id;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_users_system_contexts;
ALTER TABLE users
    DROP COLUMN IF EXISTS system_context_id;

DROP TABLE IF EXISTS system_contexts;
-- +goose StatementEnd
