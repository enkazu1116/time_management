package model

import (
	"errors"
	"time"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrWeeklyTaskIDRequired       = errors.New(messages.WeeklyTaskIDRequired)
	ErrDailyTaskIDRequired        = errors.New(messages.DailyTaskIDRequired)
	ErrDailyLogIDRequired         = errors.New(messages.DailyLogIDRequired)
	ErrWeeklyLogIDRequired        = errors.New(messages.WeeklyLogIDRequired)
	ErrUserIDRequired             = errors.New(messages.UserIDRequired)
	ErrWorkTypeRequired           = errors.New(messages.WorkTypeRequired)
	ErrWorkDateRequired           = errors.New(messages.WorkDateRequired)
	ErrWeekStartDateRequired      = errors.New(messages.WeekStartDateRequired)
	ErrTargetDurationNegative     = errors.New(messages.TargetDurationNegative)
	ErrTotalWorkDurationNegative  = errors.New(messages.TotalWorkDurationNegative)
	ErrActualDurationNegative     = errors.New(messages.ActualDurationNegative)
	ErrProgressRateOutOfRange     = errors.New(messages.ProgressRateOutOfRange)
	ErrAchievementRateOutOfRange  = errors.New(messages.AchievementRateOutOfRange)
	ErrAverageProgressOutOfRange  = errors.New(messages.AverageProgressOutOfRange)
	ErrWeeklyLogUpdatedAtRequired = errors.New(messages.WeeklyLogUpdatedAtRequired)
)

type WeeklyTask struct {
	WeeklyTaskID    string
	UserID          string
	WorkTypeID      string
	WeekStartDate   time.Time
	TargetDuration  time.Duration
	GoalDescription string
}

type DailyTask struct {
	DailyTaskID         string
	WeeklyTaskID        string
	UserID              string
	WorkDate            time.Time
	ProgressDescription string
	ProgressRate        float64
}

type DailyLog struct {
	DailyLogID        string
	UserID            string
	WorkDate          time.Time
	TotalWorkDuration time.Duration
	DailyTasks        []DailyTask
}

type WeeklyLog struct {
	WeeklyLogID         string
	WeeklyTaskID        string
	UserID              string
	WorkTypeID          string
	WeekStartDate       time.Time
	TargetDuration      time.Duration
	ActualDuration      time.Duration
	AchievementRate     float64
	GoalDescription     string
	ProgressEvaluation  string
	AverageProgressRate float64
	UpdatedAt           time.Time
}

func NewWeeklyTask(weeklyTaskID string, userID string, workTypeID string, weekStartDate time.Time, targetDuration time.Duration, goalDescription string) WeeklyTask {
	return WeeklyTask{WeeklyTaskID: weeklyTaskID, UserID: userID, WorkTypeID: workTypeID, WeekStartDate: weekStartDate, TargetDuration: targetDuration, GoalDescription: goalDescription}
}

func NewDailyTask(dailyTaskID string, weeklyTaskID string, userID string, workDate time.Time, progressDescription string, progressRate float64) DailyTask {
	return DailyTask{DailyTaskID: dailyTaskID, WeeklyTaskID: weeklyTaskID, UserID: userID, WorkDate: workDate, ProgressDescription: progressDescription, ProgressRate: progressRate}
}

func NewDailyLog(dailyLogID string, userID string, workDate time.Time, totalWorkDuration time.Duration, dailyTasks []DailyTask) DailyLog {
	return DailyLog{DailyLogID: dailyLogID, UserID: userID, WorkDate: workDate, TotalWorkDuration: totalWorkDuration, DailyTasks: dailyTasks}
}

func NewWeeklyLog(weeklyLogID string, weeklyTaskID string, userID string, workTypeID string, weekStartDate time.Time, targetDuration time.Duration, actualDuration time.Duration, achievementRate float64, goalDescription string, progressEvaluation string, averageProgressRate float64, updatedAt time.Time) WeeklyLog {
	return WeeklyLog{WeeklyLogID: weeklyLogID, WeeklyTaskID: weeklyTaskID, UserID: userID, WorkTypeID: workTypeID, WeekStartDate: weekStartDate, TargetDuration: targetDuration, ActualDuration: actualDuration, AchievementRate: achievementRate, GoalDescription: goalDescription, ProgressEvaluation: progressEvaluation, AverageProgressRate: averageProgressRate, UpdatedAt: updatedAt}
}

func (w WeeklyTask) Validate() error {
	if w.WeeklyTaskID == "" {
		return ErrWeeklyTaskIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WeeklyTaskID); err != nil {
		return err
	}
	if w.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.UserID); err != nil {
		return err
	}
	if w.WorkTypeID == "" {
		return ErrWorkTypeRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WorkTypeID); err != nil {
		return err
	}
	if w.WeekStartDate.IsZero() {
		return ErrWeekStartDateRequired
	}
	if w.TargetDuration < 0 {
		return ErrTargetDurationNegative
	}
	return nil
}

func (d DailyTask) Validate() error {
	if d.DailyTaskID == "" {
		return ErrDailyTaskIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(d.DailyTaskID); err != nil {
		return err
	}
	if d.WeeklyTaskID == "" {
		return ErrWeeklyTaskIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(d.WeeklyTaskID); err != nil {
		return err
	}
	if d.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(d.UserID); err != nil {
		return err
	}
	if d.WorkDate.IsZero() {
		return ErrWorkDateRequired
	}
	if d.ProgressRate < 0 || d.ProgressRate > 1 {
		return ErrProgressRateOutOfRange
	}
	return nil
}

func (d DailyLog) Validate() error {
	if d.DailyLogID == "" {
		return ErrDailyLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(d.DailyLogID); err != nil {
		return err
	}
	if d.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(d.UserID); err != nil {
		return err
	}
	if d.WorkDate.IsZero() {
		return ErrWorkDateRequired
	}
	if d.TotalWorkDuration < 0 {
		return ErrTotalWorkDurationNegative
	}
	return nil
}

func (w WeeklyLog) Validate() error {
	if w.WeeklyLogID == "" {
		return ErrWeeklyLogIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WeeklyLogID); err != nil {
		return err
	}
	if w.WeeklyTaskID == "" {
		return ErrWeeklyTaskIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WeeklyTaskID); err != nil {
		return err
	}
	if w.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.UserID); err != nil {
		return err
	}
	if w.WorkTypeID == "" {
		return ErrWorkTypeRequired
	}
	if err := sharedmodel.ValidateUUIDString(w.WorkTypeID); err != nil {
		return err
	}
	if w.WeekStartDate.IsZero() {
		return ErrWeekStartDateRequired
	}
	if w.TargetDuration < 0 {
		return ErrTargetDurationNegative
	}
	if w.ActualDuration < 0 {
		return ErrActualDurationNegative
	}
	if w.AchievementRate < 0 {
		return ErrAchievementRateOutOfRange
	}
	if w.AverageProgressRate < 0 || w.AverageProgressRate > 1 {
		return ErrAverageProgressOutOfRange
	}
	if w.UpdatedAt.IsZero() {
		return ErrWeeklyLogUpdatedAtRequired
	}
	return nil
}
