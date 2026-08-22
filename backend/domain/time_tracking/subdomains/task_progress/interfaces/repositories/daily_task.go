package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/task_progress/model"
)

type DailyTaskRepository interface {
	Create(dailyTask model.DailyTask) error
	Update(dailyTask model.DailyTask) error
	Delete(dailyTaskID string) error
	FindByDailyTaskID(dailyTaskID string) (model.DailyTask, error)
	FindByUserIDAndDate(userID string, workDate time.Time) ([]model.DailyTask, error)
}
