
-- +migrate Up
CREATE TABLE users (
                       id            BIGSERIAL PRIMARY KEY,
                       email         VARCHAR(255) NOT NULL UNIQUE,
                       api_key       VARCHAR(255) NOT NULL UNIQUE,
                       token_balance BIGINT NOT NULL DEFAULT 0,
                       created_at    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- +migrate Down
DROP TABLE users;