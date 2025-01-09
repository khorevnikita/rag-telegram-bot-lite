-- +goose Up

-- Создание таблицы questions
CREATE TABLE questions
(
    id                       UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    is_published             BOOLEAN                     DEFAULT false,
    is_required              BOOLEAN                     DEFAULT false,
    text                     TEXT    NOT NULL,
    hint                     TEXT,
    type                     VARCHAR NOT NULL,
    selectable_options_count INT,
    "order"                  INT     NOT NULL,
    created_at               TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at               TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at               TIMESTAMP WITHOUT TIME ZONE
);

-- Создание таблицы question_options
CREATE TABLE question_options
(
    id                      UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    question_id             UUID    NOT NULL,
    text                    VARCHAR NOT NULL,
    require_additional_text BOOLEAN                     DEFAULT false,
    created_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at              TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at              TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (question_id) REFERENCES questions (id)
);

-- Создание таблицы answers
CREATE TABLE answers
(
    id          UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    user_id     UUID,
    question_id UUID NOT NULL,
    text        TEXT,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at  TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (question_id) REFERENCES questions (id)
);

-- Создание таблицы answers
CREATE TABLE answer_options
(
    id                 UUID PRIMARY KEY            DEFAULT uuid_generate_v4(),
    user_id            UUID,
    question_id        UUID NOT NULL,
    question_option_id UUID NOT NULL,
    answer_id          UUID NOT NULL,
    created_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at         TIMESTAMP WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    deleted_at         TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (question_id) REFERENCES questions (id),
    FOREIGN KEY (question_option_id) REFERENCES question_options (id),
    FOREIGN KEY (answer_id) REFERENCES answers (id)
);

ALTER TABLE users
    ADD COLUMN state             varchar,
    ADD COLUMN state_id          UUID,
    ADD COLUMN form_completed_at TIMESTAMP;

-- +goose Down

ALTER TABLE users
    DROP COLUMN state,
    DROP COLUMN state_id,
    DROP COLUMN form_completed_at;

DROP TABLE IF EXISTS answer_options;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS question_options;
DROP TABLE IF EXISTS questions;
