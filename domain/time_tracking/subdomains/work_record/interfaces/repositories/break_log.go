package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/work_record/model"
)

type BreakLogRepository interface {
	Create(breakLog model.BreakLog) error
	Update(breakLog model.BreakLog) error
	Delete(breakLogID string) error
	FindByBreakLogID(breakLogID string) (model.BreakLog, error)
	FindByUserIDAndDate(userID string, workDate time.Time) ([]model.BreakLog, error)
	FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]model.BreakLog, error)
	FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]model.BreakLog, error)
}
