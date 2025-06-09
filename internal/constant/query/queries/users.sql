-- name: CreateUser :one
INSERT INTO users (username, email, password, profile_picture)
VALUES ($1, $2, $3, $4)
RETURNING user_id, username, email, password, profile_picture, created_at, modified_at;


-- name: GetUserByID :one
SELECT user_id, username, email, password, profile_picture, created_at, modified_at
FROM users
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT user_id, username, email, password, profile_picture, created_at, modified_at
FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT user_id, username, email, password, profile_picture, created_at, modified_at
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    username = COALESCE(sqlc.narg('username'), username), 
    email = COALESCE(sqlc.narg('email'), email), 
    password = COALESCE(sqlc.narg('password'), password),
    profile_picture = COALESCE(sqlc.narg('profile_picture'), profile_picture),

    modified_at = now()
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: ListUsers :many
SELECT user_id, username, email, password, profile_picture, created_at, modified_at
FROM users
ORDER BY created_at DESC;

-- name: SearchUser :many
SELECT * FROM users WHERE username ILIKE '%' || $1 || '%';