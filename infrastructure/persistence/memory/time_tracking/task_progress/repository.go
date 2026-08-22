package taskprogress

import (
	"time"

	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	"time_management/infrastructure/persistence/memory"
)

type WeeklyTaskRepository struct {
	store *memory.Store
}

func NewWeeklyTaskRepository(store *memory.Store) *WeeklyTaskRepository {
	return &WeeklyTaskRepository{store: store}
}

func (r *WeeklyTaskRepository) Create(weeklyTask taskprogressmodel.WeeklyTask) error {
	return r.store.CreateWeeklyTask(weeklyTask)
}

func (r *WeeklyTaskRepository) Update(weeklyTask taskprogressmodel.WeeklyTask) error {
	return r.store.UpdateWeeklyTask(weeklyTask)
}

func (r *WeeklyTaskRepository) Delete(weeklyTaskID string) error {
	return r.store.DeleteWeeklyTask(weeklyTaskID)
}

func (r *WeeklyTaskRepository) FindByWeeklyTaskID(weeklyTaskID string) (taskprogressmodel.WeeklyTask, error) {
	return r.store.FindByWeeklyTaskID(weeklyTaskID)
}

func (r *WeeklyTaskRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.WeeklyTask, error) {
	return r.store.FindWeeklyTasksByUserIDAndWeek(userID, weekStartDate)
}

type DailyTaskRepository struct {
	store *memory.Store
}

func NewDailyTaskRepository(store *memory.Store) *DailyTaskRepository {
	return &DailyTaskRepository{store: store}
}

func (r *DailyTaskRepository) Create(dailyTask taskprogressmodel.DailyTask) error {
	return r.store.CreateDailyTask(dailyTask)
}

func (r *DailyTaskRepository) Update(dailyTask taskprogressmodel.DailyTask) error {
	return r.store.UpdateDailyTask(dailyTask)
}

func (r *DailyTaskRepository) Delete(dailyTaskID string) error {
	return r.store.DeleteDailyTask(dailyTaskID)
}

func (r *DailyTaskRepository) FindByDailyTaskID(dailyTaskID string) (taskprogressmodel.DailyTask, error) {
	return r.store.FindByDailyTaskID(dailyTaskID)
}

func (r *DailyTaskRepository) FindByUserIDAndDate(userID string, workDate time.Time) ([]taskprogressmodel.DailyTask, error) {
	return r.store.FindDailyTasksByUserIDAndDate(userID, workDate)
}

type DailyLogRepository struct {
	store *memory.Store
}

func NewDailyLogRepository(store *memory.Store) *DailyLogRepository {
	return &DailyLogRepository{store: store}
}

func (r *DailyLogRepository) Create(dailyLog taskprogressmodel.DailyLog) error {
	return r.store.CreateDailyLog(dailyLog)
}

func (r *DailyLogRepository) Update(dailyLog taskprogressmodel.DailyLog) error {
	return r.store.UpdateDailyLog(dailyLog)
}

func (r *DailyLogRepository) Delete(dailyLogID string) error {
	return r.store.DeleteDailyLog(dailyLogID)
}

func (r *DailyLogRepository) FindByDailyLogID(dailyLogID string) (taskprogressmodel.DailyLog, error) {
	return r.store.FindByDailyLogID(dailyLogID)
}

func (r *DailyLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.DailyLog, error) {
	return r.store.FindDailyLogsByUserIDAndWeek(userID, weekStartDate)
}

type WeeklyLogRepository struct {
	store *memory.Store
}

func NewWeeklyLogRepository(store *memory.Store) *WeeklyLogRepository {
	return &WeeklyLogRepository{store: store}
}

func (r *WeeklyLogRepository) Create(weeklyLog taskprogressmodel.WeeklyLog) error {
	return r.store.CreateWeeklyLog(weeklyLog)
}

func (r *WeeklyLogRepository) Update(weeklyLog taskprogressmodel.WeeklyLog) error {
	return r.store.UpdateWeeklyLog(weeklyLog)
}

func (r *WeeklyLogRepository) Delete(weeklyLogID string) error {
	return r.store.DeleteWeeklyLog(weeklyLogID)
}

func (r *WeeklyLogRepository) SaveAll(weeklyLogs []taskprogressmodel.WeeklyLog) error {
	return r.store.SaveAll(weeklyLogs)
}

func (r *WeeklyLogRepository) FindByWeeklyLogID(weeklyLogID string) (taskprogressmodel.WeeklyLog, error) {
	return r.store.FindByWeeklyLogID(weeklyLogID)
}

func (r *WeeklyLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]taskprogressmodel.WeeklyLog, error) {
	return r.store.FindWeeklyLogsByUserIDAndMonth(userID, yearMonth)
}
