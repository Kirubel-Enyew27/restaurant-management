-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING user_id, username, email, password, created_at, modified_at;


-- name: GetUserByID :one
SELECT user_id, username, email, password, created_at, modified_at
FROM users
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT user_id, username, email, password, created_at, modified_at
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT user_id, username, email, password, created_at, modified_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    username = COALESCE($2, username), 
    email = COALESCE($3, email), 
    password = COALESCE($4, password),
    created_at = COALESCE($5, created_at), 
    modified_at = COALESCE($6, modified_at)
WHERE user_id = $1
RETURNING user_id, username, email, password, created_at, modified_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT user_id, username, email, password, created_at, modified_at
FROM users
ORDER BY created_at DESC;
