-- name: CreateDailyLog :one
INSERT INTO daily_logs (
  daily_log_id,
  user_id,
  work_date,
  total_work_duration
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateDailyLog :one
UPDATE daily_logs
SET
  user_id = $2,
  work_date = $3,
  total_work_duration = $4,
  updated_at = now()
WHERE daily_log_id = $1
RETURNING *;

-- name: DeleteDailyLog :exec
DELETE FROM daily_logs
WHERE daily_log_id = $1;

-- name: FindDailyLogByID :one
SELECT *
FROM daily_logs
WHERE daily_log_id = $1;

-- name: FindDailyLogsByUserIDAndWeek :many
SELECT *
FROM daily_logs
WHERE user_id = $1
  AND work_date >= $2
  AND work_date < $2::date + INTERVAL '7 days'
ORDER BY work_date;
