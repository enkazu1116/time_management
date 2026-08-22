package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/task_progress/model"
)

type WeeklyLogRepository interface {
	Create(weeklyLog model.WeeklyLog) error
	Update(weeklyLog model.WeeklyLog) error
	Delete(weeklyLogID string) error
	SaveAll(weeklyLogs []model.WeeklyLog) error
	FindByWeeklyLogID(weeklyLogID string) (model.WeeklyLog, error)
	FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]model.WeeklyLog, error)
}
