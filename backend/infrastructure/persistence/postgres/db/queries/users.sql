-- name: CreateUser :one
INSERT INTO users (
  user_id,
  password,
  name,
  email,
  regular_start,
  regular_end,
  role_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  password = $2,
  name = $3,
  email = $4,
  regular_start = $5,
  regular_end = $6,
  role_id = $7,
  updated_at = now()
WHERE user_id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;

-- name: FindUserByID :one
SELECT *
FROM users
WHERE user_id = $1;
