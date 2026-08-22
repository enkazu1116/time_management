-- name: CreateDailyTask :one
INSERT INTO daily_tasks (
  daily_task_id,
  weekly_task_id,
  user_id,
  work_date,
  progress_description,
  progress_rate
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateDailyTask :one
UPDATE daily_tasks
SET
  weekly_task_id = $2,
  user_id = $3,
  work_date = $4,
  progress_description = $5,
  progress_rate = $6,
  updated_at = now()
WHERE daily_task_id = $1
RETURNING *;

-- name: DeleteDailyTask :exec
DELETE FROM daily_tasks
WHERE daily_task_id = $1;

-- name: FindDailyTaskByID :one
SELECT *
FROM daily_tasks
WHERE daily_task_id = $1;

-- name: FindDailyTasksByUserIDAndDate :many
SELECT *
FROM daily_tasks
WHERE user_id = $1
  AND work_date = $2
ORDER BY daily_task_id;
