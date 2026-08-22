package requests

import (
	"errors"

	sharedmodel "time_management/domain/shared/model"
	"time_management/infrastructure/util/messages"
	"time_management/infrastructure/util/parser"
)

var (
	ErrUserIDRequired              = errors.New(messages.RequestUserIDRequired)
	ErrWorkTypeIDRequired          = errors.New(messages.RequestWorkTypeIDRequired)
	ErrWorkTypeNameRequired        = errors.New(messages.RequestWorkTypeNameRequired)
	ErrWorkLogIDRequired           = errors.New(messages.RequestWorkLogIDRequired)
	ErrWorkDateRequired            = errors.New(messages.RequestWorkDateRequired)
	ErrStartedAtRequired           = errors.New(messages.RequestStartedAtRequired)
	ErrEndedAtRequired             = errors.New(messages.RequestEndedAtRequired)
	ErrWeekStartDateRequired       = errors.New(messages.RequestWeekStartDateRequired)
	ErrYearMonthRequired           = errors.New(messages.YearMonthRequired)
	ErrWeeklyTaskIDRequired        = errors.New(messages.RequestWeeklyTaskIDRequired)
	ErrNegativeDuration            = errors.New(messages.RequestNegativeDuration)
	ErrProgressRateOutOfRange      = errors.New(messages.RequestProgressRateOutOfRange)
	ErrResumeDurationWithoutResume = errors.New(messages.RequestResumeDurationWithoutResume)
	ErrUnsupportedMonthlySource    = errors.New(messages.RequestUnsupportedMonthlySource)
)

type CreateWorkLogRequest struct {
	WorkLogID           string `json:"work_log_id"`
	UserID              string `json:"user_id"`
	WorkTypeID          string `json:"work_type_id"`
	WorkTypeName        string `json:"work_type_name"`
	WorkDate            string `json:"work_date"`
	StartedAt           string `json:"started_at"`
	EndedAt             string `json:"ended_at"`
	WorkDurationSeconds int64  `json:"work_duration_seconds"`
}

type CreateBreakLogRequest struct {
	BreakLogID            string `json:"break_log_id"`
	WorkLogID             string `json:"work_log_id"`
	UserID                string `json:"user_id"`
	StartedAt             string `json:"started_at"`
	EndedAt               string `json:"ended_at"`
	BreakDurationSeconds  int64  `json:"break_duration_seconds"`
	ResumedAt             string `json:"resumed_at"`
	ResumeDurationSeconds *int64 `json:"resume_duration_seconds"`
}

type CreateWeeklyTaskRequest struct {
	WeeklyTaskID          string `json:"weekly_task_id"`
	UserID                string `json:"user_id"`
	WorkTypeID            string `json:"work_type_id"`
	WeekStartDate         string `json:"week_start_date"`
	TargetDurationSeconds int64  `json:"target_duration_seconds"`
	GoalDescription       string `json:"goal_description"`
}

type CreateDailyTaskRequest struct {
	DailyTaskID         string  `json:"daily_task_id"`
	WeeklyTaskID        string  `json:"weekly_task_id"`
	UserID              string  `json:"user_id"`
	WorkDate            string  `json:"work_date"`
	ProgressDescription string  `json:"progress_description"`
	ProgressRate        float64 `json:"progress_rate"`
}

type CreateDailyLogRequest struct {
	DailyLogID               string `json:"daily_log_id"`
	UserID                   string `json:"user_id"`
	WorkDate                 string `json:"work_date"`
	TotalWorkDurationSeconds int64  `json:"total_work_duration_seconds"`
}

type CreateDailySummaryRequest struct {
	DailyAchievementID string `json:"daily_achievement_id"`
	UserID             string `json:"user_id"`
	WorkDate           string `json:"work_date"`
}

type CreateWeeklyLogRequest struct {
	WeeklyLogIDPrefix string `json:"weekly_log_id_prefix"`
	UserID            string `json:"user_id"`
	WeekStartDate     string `json:"week_start_date"`
}

type CreateMonthlySummaryRequest struct {
	MonthlyAchievementID string `json:"monthly_achievement_id"`
	UserID               string `json:"user_id"`
	YearMonth            string `json:"year_month"`
	Source               string `json:"source"`
}

func (r CreateWorkLogRequest) Validate() error {
	if r.WorkLogID != "" {
		if err := sharedmodel.ValidateUUIDString(r.WorkLogID); err != nil {
			return err
		}
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WorkTypeID == "" {
		return ErrWorkTypeIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.WorkTypeID); err != nil {
		return err
	}
	if r.WorkTypeName == "" {
		return ErrWorkTypeNameRequired
	}
	if r.WorkDate == "" {
		return ErrWorkDateRequired
	}
	if _, err := parser.ParseDate(r.WorkDate); err != nil {
		return err
	}
	if r.StartedAt == "" {
		return ErrStartedAtRequired
	}
	if _, err := parser.ParseDateTime(r.StartedAt); err != nil {
		return err
	}
	if r.EndedAt == "" {
		return ErrEndedAtRequired
	}
	if _, err := parser.ParseDateTime(r.EndedAt); err != nil {
		return err
	}
	if r.WorkDurationSeconds < 0 {
		return ErrNegativeDuration
	}
	return nil
}

func (r CreateBreakLogRequest) Validate() error {
	if r.BreakLogID != "" {
		if err := sharedmodel.ValidateUUIDString(r.BreakLogID); err != nil {
			return err
		}
	}
	if r.WorkLogID == "" {
		return ErrWorkLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.WorkLogID); err != nil {
		return err
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.StartedAt == "" {
		return ErrStartedAtRequired
	}
	if _, err := parser.ParseDateTime(r.StartedAt); err != nil {
		return err
	}
	if r.EndedAt == "" {
		return ErrEndedAtRequired
	}
	if _, err := parser.ParseDateTime(r.EndedAt); err != nil {
		return err
	}
	if r.BreakDurationSeconds < 0 {
		return ErrNegativeDuration
	}
	if r.ResumedAt != "" {
		if _, err := parser.ParseDateTime(r.ResumedAt); err != nil {
			return err
		}
	}
	if r.ResumeDurationSeconds != nil && *r.ResumeDurationSeconds < 0 {
		return ErrNegativeDuration
	}
	if r.ResumeDurationSeconds != nil && r.ResumedAt == "" {
		return ErrResumeDurationWithoutResume
	}
	return nil
}

func (r CreateWeeklyTaskRequest) Validate() error {
	if r.WeeklyTaskID != "" {
		if err := sharedmodel.ValidateUUIDString(r.WeeklyTaskID); err != nil {
			return err
		}
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WorkTypeID == "" {
		return ErrWorkTypeIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.WorkTypeID); err != nil {
		return err
	}
	if r.WeekStartDate == "" {
		return ErrWeekStartDateRequired
	}
	if _, err := parser.ParseDate(r.WeekStartDate); err != nil {
		return err
	}
	if r.TargetDurationSeconds < 0 {
		return ErrNegativeDuration
	}
	return nil
}

func (r CreateDailyTaskRequest) Validate() error {
	if r.DailyTaskID != "" {
		if err := sharedmodel.ValidateUUIDString(r.DailyTaskID); err != nil {
			return err
		}
	}
	if r.WeeklyTaskID == "" {
		return ErrWeeklyTaskIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.WeeklyTaskID); err != nil {
		return err
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WorkDate == "" {
		return ErrWorkDateRequired
	}
	if _, err := parser.ParseDate(r.WorkDate); err != nil {
		return err
	}
	if r.ProgressRate < 0 || r.ProgressRate > 1 {
		return ErrProgressRateOutOfRange
	}
	return nil
}

func (r CreateDailyLogRequest) Validate() error {
	if r.DailyLogID != "" {
		if err := sharedmodel.ValidateUUIDString(r.DailyLogID); err != nil {
			return err
		}
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WorkDate == "" {
		return ErrWorkDateRequired
	}
	if _, err := parser.ParseDate(r.WorkDate); err != nil {
		return err
	}
	if r.TotalWorkDurationSeconds < 0 {
		return ErrNegativeDuration
	}
	return nil
}

func (r CreateDailySummaryRequest) Validate() error {
	if r.DailyAchievementID != "" {
		if err := sharedmodel.ValidateUUIDString(r.DailyAchievementID); err != nil {
			return err
		}
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WorkDate == "" {
		return ErrWorkDateRequired
	}
	_, err := parser.ParseDate(r.WorkDate)
	return err
}

func (r CreateWeeklyLogRequest) Validate() error {
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.WeekStartDate == "" {
		return ErrWeekStartDateRequired
	}
	_, err := parser.ParseDate(r.WeekStartDate)
	return err
}

func (r CreateMonthlySummaryRequest) Validate() error {
	if r.MonthlyAchievementID != "" {
		if err := sharedmodel.ValidateUUIDString(r.MonthlyAchievementID); err != nil {
			return err
		}
	}
	if r.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(r.UserID); err != nil {
		return err
	}
	if r.YearMonth == "" {
		return ErrYearMonthRequired
	}
	if _, err := parser.ParseYearMonth(r.YearMonth); err != nil {
		return err
	}
	if r.Source != "" && r.Source != "work_logs" {
		return ErrUnsupportedMonthlySource
	}
	return nil
}
