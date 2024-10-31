-- +goose Up
create table USERS (
    id         uuid primary key default gen_random_uuid(),
    name       varchar not null unique,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

-- +goose Down
drop table USERS;