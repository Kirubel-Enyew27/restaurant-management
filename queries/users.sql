-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING user_id, username, email, password, created_at;


-- name: GetUserByID :one
SELECT user_id, username, email, password, created_at
FROM users
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT user_id, username, email, password, created_at
FROM users
ORDER BY created_at DESC;
