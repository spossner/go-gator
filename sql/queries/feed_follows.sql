-- name: CreateFeedFollow :one
WITH ff as (
    INSERT INTO feed_follows (user_id, feed_id)
    VALUES ($1, $2)
    RETURNING *
)
SELECT
    ff.*,
    f.name as feed_name,
    u.name as user_name
from ff
join users u on ff.user_id = u.id
join feeds f on ff.feed_id = f.id;

-- name: DeleteFeedFollow :one
WITH ff as (
    DELETE FROM feed_follows
    WHERE feed_follows.user_id = $1 and feed_follows.feed_id = $2
    RETURNING *
)
SELECT
    ff.*,
    f.name as feed_name,
    u.name as user_name
from ff
         join users u on ff.user_id = u.id
         join feeds f on ff.feed_id = f.id;

-- name: GetFeedFollowsForUser :many
select
    ff.*,
    f.name as feed_name,
    u.name as user_name
from feed_follows ff
join users u on ff.user_id = u.id
join feeds f on ff.feed_id = f.id
where ff.user_id = $1;