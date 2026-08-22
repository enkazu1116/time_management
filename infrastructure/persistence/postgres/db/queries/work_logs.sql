-- name: CreateWorkLog :one
INSERT INTO work_logs (
  work_log_id,
  user_id,
  work_type_id,
  work_date,
  started_at,
  ended_at,
  work_duration
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateWorkLog :one
UPDATE work_logs
SET
  user_id = $2,
  work_type_id = $3,
  work_date = $4,
  started_at = $5,
  ended_at = $6,
  work_duration = $7,
  updated_at = now()
WHERE work_log_id = $1
RETURNING *;

-- name: DeleteWorkLog :exec
DELETE FROM work_logs
WHERE work_log_id = $1;

-- name: FindWorkLogByID :one
SELECT *
FROM work_logs
WHERE work_log_id = $1;

-- name: FindWorkLogsByUserIDAndDate :many
SELECT *
FROM work_logs
WHERE user_id = $1
  AND work_date = $2
ORDER BY started_at;

-- name: FindWorkLogsByUserIDAndWeek :many
SELECT *
FROM work_logs
WHERE user_id = $1
  AND work_date >= $2
  AND work_date < $2::date + INTERVAL '7 days'
ORDER BY started_at;

-- name: FindWorkLogsByUserIDAndMonth :many
SELECT *
FROM work_logs
WHERE user_id = $1
  AND work_date >= date_trunc('month', $2::date)::date
  AND work_date < (date_trunc('month', $2::date) + INTERVAL '1 month')::date
ORDER BY started_at;
