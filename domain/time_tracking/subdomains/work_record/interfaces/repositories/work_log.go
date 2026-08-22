package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/work_record/model"
)

type WorkLogRepository interface {
	Create(workLog model.WorkLog) error
	Update(workLog model.WorkLog) error
	Delete(workLogID string) error
	FindByWorkLogID(workLogID string) (model.WorkLog, error)
	FindByUserIDAndDate(userID string, workDate time.Time) ([]model.WorkLog, error)
	FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]model.WorkLog, error)
	FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]model.WorkLog, error)
}
