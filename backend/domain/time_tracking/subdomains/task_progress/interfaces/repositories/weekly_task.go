package repositories

import (
	"time"

	"time_management/domain/time_tracking/subdomains/task_progress/model"
)

type WeeklyTaskRepository interface {
	Create(weeklyTask model.WeeklyTask) error
	Update(weeklyTask model.WeeklyTask) error
	Delete(weeklyTaskID string) error
	FindByWeeklyTaskID(weeklyTaskID string) (model.WeeklyTask, error)
	FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]model.WeeklyTask, error)
}
