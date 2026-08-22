CREATE TABLE roles (
  role_id uuid PRIMARY KEY,
  role_name text NOT NULL UNIQUE
);

CREATE TABLE users (
  user_id uuid PRIMARY KEY,
  password text NOT NULL,
  name text NOT NULL,
  email text NOT NULL UNIQUE,
  regular_start text NOT NULL,
  regular_end text NOT NULL,
  role_id uuid NOT NULL REFERENCES roles(role_id),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE work_types (
  work_type_id uuid PRIMARY KEY,
  work_type_name text NOT NULL UNIQUE,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE work_logs (
  work_log_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_type_id uuid NOT NULL REFERENCES work_types(work_type_id),
  work_date date NOT NULL,
  started_at timestamptz NOT NULL,
  ended_at timestamptz NOT NULL,
  work_duration interval NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX work_logs_user_date_idx ON work_logs(user_id, work_date);
CREATE INDEX work_logs_user_started_at_idx ON work_logs(user_id, started_at);

CREATE TABLE break_logs (
  break_log_id uuid PRIMARY KEY,
  work_log_id uuid NOT NULL REFERENCES work_logs(work_log_id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  started_at timestamptz NOT NULL,
  ended_at timestamptz NOT NULL,
  break_duration interval NOT NULL,
  resumed_at timestamptz,
  resume_duration interval,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX break_logs_user_started_at_idx ON break_logs(user_id, started_at);
CREATE INDEX break_logs_work_log_idx ON break_logs(work_log_id);

CREATE TABLE daily_achievements (
  daily_achievement_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_date date NOT NULL,
  total_work_duration interval NOT NULL,
  total_break_duration interval NOT NULL,
  work_count integer NOT NULL,
  break_count integer NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE(user_id, work_date)
);

CREATE TABLE daily_work_type_summaries (
  daily_achievement_id uuid NOT NULL REFERENCES daily_achievements(daily_achievement_id) ON DELETE CASCADE,
  work_type_id uuid NOT NULL REFERENCES work_types(work_type_id),
  work_count integer NOT NULL,
  work_duration interval NOT NULL,
  average_work_duration interval NOT NULL,
  PRIMARY KEY(daily_achievement_id, work_type_id)
);

CREATE TABLE daily_break_summaries (
  daily_achievement_id uuid PRIMARY KEY REFERENCES daily_achievements(daily_achievement_id) ON DELETE CASCADE,
  break_count integer NOT NULL,
  break_duration interval NOT NULL,
  average_break_duration interval NOT NULL,
  average_resume_duration interval NOT NULL,
  longest_resume_duration interval NOT NULL
);

CREATE TABLE weekly_tasks (
  weekly_task_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_type_id uuid NOT NULL REFERENCES work_types(work_type_id),
  week_start_date date NOT NULL,
  target_duration interval NOT NULL,
  goal_description text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX weekly_tasks_user_week_idx ON weekly_tasks(user_id, week_start_date);

CREATE TABLE daily_tasks (
  daily_task_id uuid PRIMARY KEY,
  weekly_task_id uuid NOT NULL REFERENCES weekly_tasks(weekly_task_id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_date date NOT NULL,
  progress_description text NOT NULL,
  progress_rate double precision NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX daily_tasks_user_date_idx ON daily_tasks(user_id, work_date);
CREATE INDEX daily_tasks_weekly_task_idx ON daily_tasks(weekly_task_id);

CREATE TABLE daily_logs (
  daily_log_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_date date NOT NULL,
  total_work_duration interval NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE(user_id, work_date)
);

CREATE INDEX daily_logs_user_date_idx ON daily_logs(user_id, work_date);

CREATE TABLE weekly_logs (
  weekly_log_id uuid PRIMARY KEY,
  weekly_task_id uuid NOT NULL REFERENCES weekly_tasks(weekly_task_id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  work_type_id uuid NOT NULL REFERENCES work_types(work_type_id),
  week_start_date date NOT NULL,
  target_duration interval NOT NULL,
  actual_duration interval NOT NULL,
  achievement_rate double precision NOT NULL,
  goal_description text NOT NULL,
  progress_evaluation text NOT NULL,
  average_progress_rate double precision NOT NULL,
  updated_at timestamptz NOT NULL
);

CREATE INDEX weekly_logs_user_week_idx ON weekly_logs(user_id, week_start_date);
CREATE INDEX weekly_logs_user_month_idx ON weekly_logs(user_id, week_start_date);

CREATE TABLE monthly_achievements (
  monthly_achievement_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  year_month date NOT NULL,
  total_work_duration interval NOT NULL,
  total_break_duration interval NOT NULL,
  work_days integer NOT NULL,
  work_count integer NOT NULL,
  break_count integer NOT NULL,
  updated_at timestamptz NOT NULL,
  UNIQUE(user_id, year_month)
);

CREATE TABLE monthly_work_type_summaries (
  monthly_achievement_id uuid NOT NULL REFERENCES monthly_achievements(monthly_achievement_id) ON DELETE CASCADE,
  work_type_id uuid NOT NULL REFERENCES work_types(work_type_id),
  work_count integer NOT NULL,
  work_duration interval NOT NULL,
  average_work_duration interval NOT NULL,
  PRIMARY KEY(monthly_achievement_id, work_type_id)
);

CREATE TABLE monthly_break_summaries (
  monthly_achievement_id uuid PRIMARY KEY REFERENCES monthly_achievements(monthly_achievement_id) ON DELETE CASCADE,
  break_count integer NOT NULL,
  break_duration interval NOT NULL,
  average_break_duration interval NOT NULL,
  average_resume_duration interval NOT NULL,
  longest_resume_duration interval NOT NULL
);
