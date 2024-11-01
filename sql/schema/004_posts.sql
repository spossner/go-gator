-- +goose Up
CREATE TABLE posts (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_id         uuid NOT NULL REFERENCES feeds(ID) ON DELETE CASCADE,
    title           VARCHAR NOT NULL,
    url             VARCHAR NOT NULL UNIQUE,
    description     VARCHAR,
    published_at    TIMESTAMP,
    created_at      TIMESTAMP DEFAULT current_timestamp,
    updated_at      TIMESTAMP DEFAULT current_timestamp
);

-- +goose Down
DROP TABLE posts;