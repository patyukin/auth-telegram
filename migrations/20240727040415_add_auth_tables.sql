-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE public.users
(
    id            UUID              DEFAULT uuid_generate_v4() PRIMARY KEY,
    login         VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NULL,
    surname       VARCHAR(255) NULL,
    role          VARCHAR(255) NULL DEFAULT 'user',
    created_at    TIMESTAMP    NOT NULL,
    updated_at    TIMESTAMP    NULL
);

CREATE TABLE public.tokens
(
    token      UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id    UUID      NOT NULL REFERENCES public.users (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NULL
);

CREATE TABLE public.telegram_users
(
    id          UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id     UUID UNIQUE NOT NULL REFERENCES public.users (id) ON UPDATE CASCADE ON DELETE NO ACTION,
    tg_username VARCHAR(50) NOT NULL UNIQUE,
    tg_user_id  BIGINT UNIQUE,
    tg_chat_id  BIGINT,
    created_at  TIMESTAMP   NOT NULL,
    updated_at  TIMESTAMP   NULL
);

-- +goose Down
DROP TABLE public.tokens;
DROP TABLE public.users;
DROP TABLE public.telegram_users;
