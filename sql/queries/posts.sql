-- name: CreatePost :one
INSERT INTO posts (feed_id, title, url, description, published_at)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetPostsByUser :many
SELECT p.*, f.name, f.url
FROM posts p
JOIN feeds f ON p.feed_id = f.id
JOIN feed_follows ff on ff.feed_id = f.id
WHERE ff.user_id = $1
ORDER BY p.published_at desc
LIMIT $2
OFFSET $3;
