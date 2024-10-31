-- +goose Up
create table FEED_FOLLOWS (
    id         uuid primary key default gen_random_uuid(),
    user_id    uuid not null references users(id) on delete cascade,
    feed_id    uuid not null references feeds(id) on delete cascade,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    UNIQUE(user_id, feed_id)
);

-- +goose Down
drop table FEED_FOLLOWS;