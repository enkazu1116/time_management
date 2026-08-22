-- name: CreateWeeklyTask :one
INSERT INTO weekly_tasks (
  weekly_task_id,
  user_id,
  work_type_id,
  week_start_date,
  target_duration,
  goal_description
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateWeeklyTask :one
UPDATE weekly_tasks
SET
  user_id = $2,
  work_type_id = $3,
  week_start_date = $4,
  target_duration = $5,
  goal_description = $6,
  updated_at = now()
WHERE weekly_task_id = $1
RETURNING *;

-- name: DeleteWeeklyTask :exec
DELETE FROM weekly_tasks
WHERE weekly_task_id = $1;

-- name: FindWeeklyTaskByID :one
SELECT *
FROM weekly_tasks
WHERE weekly_task_id = $1;

-- name: FindWeeklyTasksByUserIDAndWeek :many
SELECT *
FROM weekly_tasks
WHERE user_id = $1
  AND week_start_date = $2
ORDER BY work_type_id;
