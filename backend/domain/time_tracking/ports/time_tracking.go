package ports

import "time"

type CreateWorkLogInput struct {
	WorkLogID, UserID, WorkTypeID, WorkTypeName string
	WorkDate, StartedAt, EndedAt                time.Time
	WorkDuration                                time.Duration
}
type CreateBreakLogInput struct {
	BreakLogID, WorkLogID, UserID string
	StartedAt, EndedAt            time.Time
	BreakDuration                 time.Duration
	ResumedAt                     *time.Time
	ResumeDuration                *time.Duration
}
type CreateWeeklyTaskInput struct {
	WeeklyTaskID, UserID, WorkTypeID string
	WeekStartDate                    time.Time
	TargetDuration                   time.Duration
	GoalDescription                  string
}
type CreateDailyTaskInput struct {
	DailyTaskID, WeeklyTaskID, UserID string
	WorkDate                          time.Time
	ProgressDescription               string
	ProgressRate                      float64
}
type CreateDailyLogInput struct {
	DailyLogID, UserID string
	WorkDate           time.Time
	TotalWorkDuration  time.Duration
}

type CreateWeeklyLogInput struct {
	WeeklyLogIDPrefix, UserID string
	WeekStartDate, UpdatedAt  time.Time
}
type CreateDailySummaryInput struct {
	DailyAchievementID, UserID string
	WorkDate, UpdatedAt        time.Time
}
type CreateMonthlySummaryInput struct {
	MonthlyAchievementID, UserID string
	YearMonth                    time.Time
	UpdatedAt                    time.Time
}

type WorkLogOutput struct {
	WorkLogID, UserID, WorkTypeID, WorkTypeName string
	WorkDate, StartedAt, EndedAt                time.Time
	WorkDuration                                time.Duration
}
type BreakLogOutput struct {
	BreakLogID, WorkLogID, UserID string
	StartedAt, EndedAt            time.Time
	BreakDuration                 time.Duration
	ResumedAt                     *time.Time
	ResumeDuration                *time.Duration
}
type WeeklyTaskOutput struct {
	WeeklyTaskID, UserID, WorkTypeID string
	WeekStartDate                    time.Time
	TargetDuration                   time.Duration
	GoalDescription                  string
}
type DailyTaskOutput struct {
	DailyTaskID, WeeklyTaskID, UserID string
	WorkDate                          time.Time
	ProgressDescription               string
	ProgressRate                      float64
}
type DailyLogOutput struct {
	DailyLogID, UserID string
	WorkDate           time.Time
	TotalWorkDuration  time.Duration
}

type WeeklyLogOutput struct {
	WeeklyLogID, WeeklyTaskID, UserID, WorkTypeID string
	WeekStartDate                                 time.Time
	TargetDuration, ActualDuration                time.Duration
	AchievementRate                               float64
	ProgressEvaluation                            string
}
type WeeklyLogSetOutput struct{ WeeklyLogs []WeeklyLogOutput }
type SummaryOutput struct {
	Kind                                  string
	AchievementID, UserID, Period         string
	TotalWorkDuration, TotalBreakDuration time.Duration
	WorkDays, WorkCount, BreakCount       int
}

type TimeTrackingUsecase interface {
	CreateWorkLog(CreateWorkLogInput) (WorkLogOutput, error)
	CreateBreakLog(CreateBreakLogInput) (BreakLogOutput, error)
	CreateWeeklyTask(CreateWeeklyTaskInput) (WeeklyTaskOutput, error)
	CreateDailyTask(CreateDailyTaskInput) (DailyTaskOutput, error)
	CreateDailyLog(CreateDailyLogInput) (DailyLogOutput, error)
	CreateDailySummary(CreateDailySummaryInput) (SummaryOutput, error)
	CreateWeeklyLog(CreateWeeklyLogInput) (WeeklyLogSetOutput, error)
	CreateMonthlySummary(CreateMonthlySummaryInput, bool) (SummaryOutput, error)
}
