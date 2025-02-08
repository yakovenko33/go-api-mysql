-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id char(36) not null primary key,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    email varchar(150) not null,
    password varchar(255) not null,
    created_at timestamp null,
    updated_at timestamp null,
    deleted_at timestamp null,
    status tinyint(1) default 0 not null,
    created_by varchar(36) null,
    modified_user_id varchar(36) null,
    constraint users_email_unique unique (email)
);

CREATE INDEX users_created_by_index ON users (created_by);
CREATE INDEX users_modified_user_id_index ON users (modified_user_id);

-- +goose Down
DROP INDEX users_created_by_index;
DROP INDEX users_modified_user_id_index;
DROP TABLE IF EXISTS users;

