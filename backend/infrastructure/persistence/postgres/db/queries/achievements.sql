-- name: UpsertDailyAchievement :one
INSERT INTO daily_achievements (
  daily_achievement_id,
  user_id,
  work_date,
  total_work_duration,
  total_break_duration,
  work_count,
  break_count,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (user_id, work_date)
DO UPDATE SET
  daily_achievement_id = EXCLUDED.daily_achievement_id,
  total_work_duration = EXCLUDED.total_work_duration,
  total_break_duration = EXCLUDED.total_break_duration,
  work_count = EXCLUDED.work_count,
  break_count = EXCLUDED.break_count,
  updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: DeleteDailyAchievement :exec
DELETE FROM daily_achievements
WHERE daily_achievement_id = $1;

-- name: FindDailyAchievementByUserIDAndDate :one
SELECT *
FROM daily_achievements
WHERE user_id = $1
  AND work_date = $2;

-- name: UpsertDailyWorkTypeSummary :one
INSERT INTO daily_work_type_summaries (
  daily_achievement_id,
  work_type_id,
  work_count,
  work_duration,
  average_work_duration
) VALUES (
  $1, $2, $3, $4, $5
)
ON CONFLICT (daily_achievement_id, work_type_id)
DO UPDATE SET
  work_count = EXCLUDED.work_count,
  work_duration = EXCLUDED.work_duration,
  average_work_duration = EXCLUDED.average_work_duration
RETURNING *;

-- name: DeleteDailyWorkTypeSummaries :exec
DELETE FROM daily_work_type_summaries
WHERE daily_achievement_id = $1;

-- name: FindDailyWorkTypeSummaries :many
SELECT *
FROM daily_work_type_summaries
WHERE daily_achievement_id = $1
ORDER BY work_type_id;

-- name: UpsertDailyBreakSummary :one
INSERT INTO daily_break_summaries (
  daily_achievement_id,
  break_count,
  break_duration,
  average_break_duration,
  average_resume_duration,
  longest_resume_duration
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (daily_achievement_id)
DO UPDATE SET
  break_count = EXCLUDED.break_count,
  break_duration = EXCLUDED.break_duration,
  average_break_duration = EXCLUDED.average_break_duration,
  average_resume_duration = EXCLUDED.average_resume_duration,
  longest_resume_duration = EXCLUDED.longest_resume_duration
RETURNING *;

-- name: FindDailyBreakSummary :one
SELECT *
FROM daily_break_summaries
WHERE daily_achievement_id = $1;

-- name: UpsertMonthlyAchievement :one
INSERT INTO monthly_achievements (
  monthly_achievement_id,
  user_id,
  year_month,
  total_work_duration,
  total_break_duration,
  work_days,
  work_count,
  break_count,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
ON CONFLICT (user_id, year_month)
DO UPDATE SET
  monthly_achievement_id = EXCLUDED.monthly_achievement_id,
  total_work_duration = EXCLUDED.total_work_duration,
  total_break_duration = EXCLUDED.total_break_duration,
  work_days = EXCLUDED.work_days,
  work_count = EXCLUDED.work_count,
  break_count = EXCLUDED.break_count,
  updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: DeleteMonthlyAchievement :exec
DELETE FROM monthly_achievements
WHERE monthly_achievement_id = $1;

-- name: FindMonthlyAchievementByUserIDAndMonth :one
SELECT *
FROM monthly_achievements
WHERE user_id = $1
  AND year_month = date_trunc('month', $2::date)::date;

-- name: UpsertMonthlyWorkTypeSummary :one
INSERT INTO monthly_work_type_summaries (
  monthly_achievement_id,
  work_type_id,
  work_count,
  work_duration,
  average_work_duration
) VALUES (
  $1, $2, $3, $4, $5
)
ON CONFLICT (monthly_achievement_id, work_type_id)
DO UPDATE SET
  work_count = EXCLUDED.work_count,
  work_duration = EXCLUDED.work_duration,
  average_work_duration = EXCLUDED.average_work_duration
RETURNING *;

-- name: DeleteMonthlyWorkTypeSummaries :exec
DELETE FROM monthly_work_type_summaries
WHERE monthly_achievement_id = $1;

-- name: FindMonthlyWorkTypeSummaries :many
SELECT *
FROM monthly_work_type_summaries
WHERE monthly_achievement_id = $1
ORDER BY work_type_id;

-- name: UpsertMonthlyBreakSummary :one
INSERT INTO monthly_break_summaries (
  monthly_achievement_id,
  break_count,
  break_duration,
  average_break_duration,
  average_resume_duration,
  longest_resume_duration
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (monthly_achievement_id)
DO UPDATE SET
  break_count = EXCLUDED.break_count,
  break_duration = EXCLUDED.break_duration,
  average_break_duration = EXCLUDED.average_break_duration,
  average_resume_duration = EXCLUDED.average_resume_duration,
  longest_resume_duration = EXCLUDED.longest_resume_duration
RETURNING *;

-- name: FindMonthlyBreakSummary :one
SELECT *
FROM monthly_break_summaries
WHERE monthly_achievement_id = $1;
