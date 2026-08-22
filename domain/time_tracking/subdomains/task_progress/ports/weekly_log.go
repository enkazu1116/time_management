package ports

import (
	"errors"
	"time"

	"time_management/domain/shared/messages"
	sharedmodel "time_management/domain/shared/model"
)

var (
	ErrWeeklyLogIDPrefixRequired = errors.New(messages.WeeklyLogIDPrefixRequired)
	ErrUserIDRequired            = errors.New(messages.UserIDRequired)
	ErrWeekStartDateRequired     = errors.New(messages.WeekStartDateRequired)
	ErrUpdatedAtRequired         = errors.New(messages.UpdatedAtRequired)
)

type CreateWeeklyLogInput struct {
	WeeklyLogIDPrefix string
	UserID            string
	WeekStartDate     time.Time
	UpdatedAt         time.Time
}

func (i CreateWeeklyLogInput) Validate() error {
	if i.WeeklyLogIDPrefix == "" {
		return ErrWeeklyLogIDPrefixRequired
	}
	if i.UserID == "" {
		return ErrUserIDRequired
	}
	if err := sharedmodel.ValidateUUIDString(i.UserID); err != nil {
		return err
	}
	if i.WeekStartDate.IsZero() {
		return ErrWeekStartDateRequired
	}
	if i.UpdatedAt.IsZero() {
		return ErrUpdatedAtRequired
	}
	return nil
}
