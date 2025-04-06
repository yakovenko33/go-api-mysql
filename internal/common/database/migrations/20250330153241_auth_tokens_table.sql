-- +goose Up
CREATE TABLE IF NOT EXISTS auth_tokens (
    id char(36) not null primary key,
    user_agent varchar(100) not null,
    access_token varchar(4000) not null,
    refresh_token varchar(4000) not null,
    created_at timestamp null,
    expires_in timestamp null
);

-- +goose Down
DROP TABLE IF EXISTS auth_tokens;
