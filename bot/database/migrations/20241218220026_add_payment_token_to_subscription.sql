-- +goose Up
-- +goose StatementBegin

CREATE TABLE payment_providers
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP        DEFAULT NOW(),
    updated_at TIMESTAMP        DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE payment_methods
(
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id           UUID         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider_id       UUID         NOT NULL REFERENCES payment_providers (id) ON DELETE CASCADE,
    token             VARCHAR(255) NOT NULL,
    card_last_four    VARCHAR(4),
    card_expiry_month INT,
    card_expiry_year  INT,
    created_at        TIMESTAMP        DEFAULT NOW(),
    updated_at        TIMESTAMP        DEFAULT NOW(),
    deleted_at        TIMESTAMP
);

ALTER TABLE subscriptions
    ADD COLUMN unsubscribed_at   TIMESTAMP, -- Ссылка на метод оплаты
    ADD COLUMN payment_method_id UUID, -- Ссылка на метод оплаты
    ADD COLUMN provider_id       UUID; -- Ссылка на провайдера

ALTER TABLE subscriptions
    ADD CONSTRAINT fk_payment_method FOREIGN KEY (payment_method_id) REFERENCES payment_methods (id),
    ADD CONSTRAINT fk_provider FOREIGN KEY (provider_id) REFERENCES payment_providers (id);


CREATE TABLE payments
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subscription_id UUID           NOT NULL REFERENCES subscriptions (id) ON DELETE CASCADE,
    provider_id     UUID           NOT NULL REFERENCES payment_providers (id) ON DELETE CASCADE,
    amount          DECIMAL(10, 2) NOT NULL,
    currency        VARCHAR(10)      DEFAULT 'RUB',
    created_at      TIMESTAMP        DEFAULT NOW(),
    updated_at      TIMESTAMP        DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS payments;

ALTER TABLE subscriptions
    DROP CONSTRAINT fk_payment_method;
ALTER TABLE subscriptions
    DROP CONSTRAINT fk_provider;

ALTER TABLE subscriptions
    DROP COLUMN IF EXISTS unsubscribed_at,
    DROP COLUMN IF EXISTS payment_method_id,
    DROP COLUMN IF EXISTS provider_id;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS payment_providers;

-- +goose StatementEnd
