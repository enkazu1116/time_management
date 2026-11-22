-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUserName :many
SELECT * FROM users WHERE user_name = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByUserNameAndEmail :one
SELECT * FROM users WHERE user_name = $1 AND email = $2 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, user_name, password, email, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET user_name = $2, password = $3, email = $4, updated_at = $5 WHERE id = $1 RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id = $1 RETURNING *;