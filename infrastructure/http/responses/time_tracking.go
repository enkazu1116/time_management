package responses

import (
	"time"
	trackingports "time_management/domain/time_tracking/ports"
)

type WorkLogResponse struct {
	WorkLogID           string `json:"work_log_id"`
	UserID              string `json:"user_id"`
	WorkTypeID          string `json:"work_type_id"`
	WorkTypeName        string `json:"work_type_name"`
	WorkDate            string `json:"work_date"`
	StartedAt           string `json:"started_at"`
	EndedAt             string `json:"ended_at"`
	WorkDurationSeconds int64  `json:"work_duration_seconds"`
}
type BreakLogResponse struct {
	BreakLogID            string  `json:"break_log_id"`
	WorkLogID             string  `json:"work_log_id"`
	UserID                string  `json:"user_id"`
	StartedAt             string  `json:"started_at"`
	EndedAt               string  `json:"ended_at"`
	BreakDurationSeconds  int64   `json:"break_duration_seconds"`
	ResumedAt             *string `json:"resumed_at,omitempty"`
	ResumeDurationSeconds *int64  `json:"resume_duration_seconds,omitempty"`
}
type WeeklyTaskResponse struct {
	WeeklyTaskID          string `json:"weekly_task_id"`
	UserID                string `json:"user_id"`
	WorkTypeID            string `json:"work_type_id"`
	WeekStartDate         string `json:"week_start_date"`
	TargetDurationSeconds int64  `json:"target_duration_seconds"`
	GoalDescription       string `json:"goal_description"`
}
type DailyTaskResponse struct {
	DailyTaskID         string  `json:"daily_task_id"`
	WeeklyTaskID        string  `json:"weekly_task_id"`
	UserID              string  `json:"user_id"`
	WorkDate            string  `json:"work_date"`
	ProgressDescription string  `json:"progress_description"`
	ProgressRate        float64 `json:"progress_rate"`
}
type DailyLogResponse struct {
	DailyLogID               string `json:"daily_log_id"`
	UserID                   string `json:"user_id"`
	WorkDate                 string `json:"work_date"`
	TotalWorkDurationSeconds int64  `json:"total_work_duration_seconds"`
}

func formatDate(v time.Time) string { return v.Format("2006-01-02") }
func formatTime(v time.Time) string { return v.Format(time.RFC3339) }
func WorkLog(v trackingports.WorkLogOutput) WorkLogResponse {
	return WorkLogResponse{v.WorkLogID, v.UserID, v.WorkTypeID, v.WorkTypeName, formatDate(v.WorkDate), formatTime(v.StartedAt), formatTime(v.EndedAt), int64(v.WorkDuration / time.Second)}
}
func BreakLog(v trackingports.BreakLogOutput) BreakLogResponse {
	r := BreakLogResponse{BreakLogID: v.BreakLogID, WorkLogID: v.WorkLogID, UserID: v.UserID, StartedAt: formatTime(v.StartedAt), EndedAt: formatTime(v.EndedAt), BreakDurationSeconds: int64(v.BreakDuration / time.Second)}
	if v.ResumedAt != nil {
		s := formatTime(*v.ResumedAt)
		r.ResumedAt = &s
	}
	if v.ResumeDuration != nil {
		n := int64(*v.ResumeDuration / time.Second)
		r.ResumeDurationSeconds = &n
	}
	return r
}
func WeeklyTask(v trackingports.WeeklyTaskOutput) WeeklyTaskResponse {
	return WeeklyTaskResponse{v.WeeklyTaskID, v.UserID, v.WorkTypeID, formatDate(v.WeekStartDate), int64(v.TargetDuration / time.Second), v.GoalDescription}
}
func DailyTask(v trackingports.DailyTaskOutput) DailyTaskResponse {
	return DailyTaskResponse{v.DailyTaskID, v.WeeklyTaskID, v.UserID, formatDate(v.WorkDate), v.ProgressDescription, v.ProgressRate}
}
func DailyLog(v trackingports.DailyLogOutput) DailyLogResponse {
	return DailyLogResponse{v.DailyLogID, v.UserID, formatDate(v.WorkDate), int64(v.TotalWorkDuration / time.Second)}
}

type WeeklyLogResponse struct {
	WeeklyLogID           string  `json:"weekly_log_id"`
	WeeklyTaskID          string  `json:"weekly_task_id"`
	UserID                string  `json:"user_id"`
	WorkTypeID            string  `json:"work_type_id"`
	WeekStartDate         string  `json:"week_start_date"`
	TargetDurationSeconds int64   `json:"target_duration_seconds"`
	ActualDurationSeconds int64   `json:"actual_duration_seconds"`
	AchievementRate       float64 `json:"achievement_rate"`
	ProgressEvaluation    string  `json:"progress_evaluation"`
}
type WeeklyLogSetResponse struct {
	WeeklyLogs []WeeklyLogResponse `json:"weekly_logs"`
}

func WeeklyLogSet(v trackingports.WeeklyLogSetOutput) WeeklyLogSetResponse {
	out := WeeklyLogSetResponse{}
	for _, x := range v.WeeklyLogs {
		out.WeeklyLogs = append(out.WeeklyLogs, WeeklyLogResponse{x.WeeklyLogID, x.WeeklyTaskID, x.UserID, x.WorkTypeID, formatDate(x.WeekStartDate), int64(x.TargetDuration / time.Second), int64(x.ActualDuration / time.Second), x.AchievementRate, x.ProgressEvaluation})
	}
	return out
}

type SummaryResponse struct {
	Kind                      string `json:"kind"`
	AchievementID             string `json:"achievement_id"`
	UserID                    string `json:"user_id"`
	Period                    string `json:"period"`
	TotalWorkDurationSeconds  int64  `json:"total_work_duration_seconds"`
	TotalBreakDurationSeconds int64  `json:"total_break_duration_seconds"`
	WorkDays                  int    `json:"work_days"`
	WorkCount                 int    `json:"work_count"`
	BreakCount                int    `json:"break_count"`
}

func Summary(v trackingports.SummaryOutput) SummaryResponse {
	return SummaryResponse{Kind: v.Kind, AchievementID: v.AchievementID, UserID: v.UserID, Period: v.Period, TotalWorkDurationSeconds: int64(v.TotalWorkDuration / time.Second), TotalBreakDurationSeconds: int64(v.TotalBreakDuration / time.Second), WorkDays: v.WorkDays, WorkCount: v.WorkCount, BreakCount: v.BreakCount}
}
