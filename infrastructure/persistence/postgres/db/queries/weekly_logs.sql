-- name: CreateWeeklyLog :one
INSERT INTO weekly_logs (
  weekly_log_id,
  weekly_task_id,
  user_id,
  work_type_id,
  week_start_date,
  target_duration,
  actual_duration,
  achievement_rate,
  goal_description,
  progress_evaluation,
  average_progress_rate,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateWeeklyLog :one
UPDATE weekly_logs
SET
  weekly_task_id = $2,
  user_id = $3,
  work_type_id = $4,
  week_start_date = $5,
  target_duration = $6,
  actual_duration = $7,
  achievement_rate = $8,
  goal_description = $9,
  progress_evaluation = $10,
  average_progress_rate = $11,
  updated_at = $12
WHERE weekly_log_id = $1
RETURNING *;

-- name: DeleteWeeklyLog :exec
DELETE FROM weekly_logs
WHERE weekly_log_id = $1;

-- name: FindWeeklyLogByID :one
SELECT *
FROM weekly_logs
WHERE weekly_log_id = $1;

-- name: FindWeeklyLogsByUserIDAndMonth :many
SELECT *
FROM weekly_logs
WHERE user_id = $1
  AND week_start_date >= date_trunc('month', $2::date)::date
  AND week_start_date < (date_trunc('month', $2::date) + INTERVAL '1 month')::date
ORDER BY week_start_date, work_type_id;
