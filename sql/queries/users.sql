-- name: CreateUser :one
INSERT INTO users (name)
VALUES ($1)
RETURNING *;

-- name: GetUserByName :one
SELECT *
FROM users
WHERE name = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;


-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT *
from users;