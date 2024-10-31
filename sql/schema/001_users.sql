-- +goose Up
create table USERS (
    id         uuid primary key,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    name       text not null unique
);

-- +goose Down
drop table USERS;