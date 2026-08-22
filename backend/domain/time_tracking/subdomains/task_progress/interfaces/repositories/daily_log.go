package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/task_progress/model"
)

type DailyLogRepository interface {
	Create(dailyLog model.DailyLog) error
	Update(dailyLog model.DailyLog) error
	Delete(dailyLogID string) error
	FindByDailyLogID(dailyLogID string) (model.DailyLog, error)
	FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]model.DailyLog, error)
}
