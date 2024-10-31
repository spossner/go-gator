-- +goose Up
create table FEEDS (
    id         uuid primary key default gen_random_uuid(),
    name       varchar not null,
    url        varchar not null unique,
    user_id    uuid not null references users(id) on delete cascade,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

-- +goose Down
drop table FEEDS;