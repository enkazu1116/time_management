package ports

import (
	"errors"
	"time"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrUserIDRequired               = errors.New(messages.UserIDRequired)
	ErrDailyAchievementIDRequired   = errors.New(messages.DailyAchievementIDRequired)
	ErrMonthlyAchievementIDRequired = errors.New(messages.MonthlyAchievementIDRequired)
	ErrWorkDateRequired             = errors.New(messages.WorkDateRequired)
	ErrYearMonthRequired            = errors.New(messages.YearMonthRequired)
	ErrUpdatedAtRequired            = errors.New(messages.UpdatedAtRequired)
)

type CreateDailySummaryInput struct {
	DailyAchievementID string
	UserID             string
	WorkDate           time.Time
	UpdatedAt          time.Time
}

type CreateMonthlySummaryInput struct {
	MonthlyAchievementID string
	UserID               string
	YearMonth            time.Time
	UpdatedAt            time.Time
}

func (i CreateDailySummaryInput) Validate() error {
	if i.DailyAchievementID == "" {
		return ErrDailyAchievementIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.DailyAchievementID); err != nil {
		return err
	}
	if i.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.UserID); err != nil {
		return err
	}
	if i.WorkDate.IsZero() {
		return ErrWorkDateRequired
	}
	if i.UpdatedAt.IsZero() {
		return ErrUpdatedAtRequired
	}
	return nil
}

func (i CreateMonthlySummaryInput) Validate() error {
	if i.MonthlyAchievementID == "" {
		return ErrMonthlyAchievementIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.MonthlyAchievementID); err != nil {
		return err
	}
	if i.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.UserID); err != nil {
		return err
	}
	if i.YearMonth.IsZero() {
		return ErrYearMonthRequired
	}
	if i.UpdatedAt.IsZero() {
		return ErrUpdatedAtRequired
	}
	return nil
}
