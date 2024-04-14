CREATE TABLE users
(
    id            SERIAL       PRIMARY KEY,
    nickname      VARCHAR(31)  NOT NULL,
    email         VARCHAR(63)  NOT NULL,
    password_hash varchar(255) NOT NULL,
    role          varchar(15)  NOT NULL
);

CREATE TABLE banners
(
    id            SERIAL       PRIMARY KEY,
    tag_ids       INTEGER[]    NOT NULL,
    feature_id    INTEGER      NOT NULL,
    title         VARCHAR(255) NOT NULL,
    text          TEXT         NOT NULL,
    url           varchar(255) NOT NULL,
    is_active     BOOLEAN      NOT NULL,
    created_at    TIMESTAMP    NOT NULL,
    updated_at    TIMESTAMP    NOT NULL
);