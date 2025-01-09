-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions
(
    id         UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    user_id    UUID                        NOT NULL,
    amount     FLOAT                       NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
