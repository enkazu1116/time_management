-- name: CreateBreakLog :one
INSERT INTO break_logs (
  break_log_id,
  work_log_id,
  user_id,
  started_at,
  ended_at,
  break_duration,
  resumed_at,
  resume_duration
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateBreakLog :one
UPDATE break_logs
SET
  work_log_id = $2,
  user_id = $3,
  started_at = $4,
  ended_at = $5,
  break_duration = $6,
  resumed_at = $7,
  resume_duration = $8,
  updated_at = now()
WHERE break_log_id = $1
RETURNING *;

-- name: DeleteBreakLog :exec
DELETE FROM break_logs
WHERE break_log_id = $1;

-- name: FindBreakLogByID :one
SELECT *
FROM break_logs
WHERE break_log_id = $1;

-- name: FindBreakLogsByUserIDAndDate :many
SELECT *
FROM break_logs
WHERE user_id = $1
  AND started_at >= $2::date
  AND started_at < $2::date + INTERVAL '1 day'
ORDER BY started_at;

-- name: FindBreakLogsByUserIDAndWeek :many
SELECT *
FROM break_logs
WHERE user_id = $1
  AND started_at >= $2::date
  AND started_at < $2::date + INTERVAL '7 days'
ORDER BY started_at;

-- name: FindBreakLogsByUserIDAndMonth :many
SELECT *
FROM break_logs
WHERE user_id = $1
  AND started_at >= date_trunc('month', $2::date)
  AND started_at < date_trunc('month', $2::date) + INTERVAL '1 month'
ORDER BY started_at;
