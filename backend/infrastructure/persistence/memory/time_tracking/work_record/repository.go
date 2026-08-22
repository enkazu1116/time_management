package workrecord

import (
	"time"

	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
	"time_management/infrastructure/persistence/memory"
)

type WorkLogRepository struct {
	store *memory.Store
}

func NewWorkLogRepository(store *memory.Store) *WorkLogRepository {
	return &WorkLogRepository{store: store}
}

func (r *WorkLogRepository) Create(workLog workrecordmodel.WorkLog) error {
	return r.store.CreateWorkLog(workLog)
}

func (r *WorkLogRepository) Update(workLog workrecordmodel.WorkLog) error {
	return r.store.UpdateWorkLog(workLog)
}

func (r *WorkLogRepository) Delete(workLogID string) error {
	return r.store.DeleteWorkLog(workLogID)
}

func (r *WorkLogRepository) FindByWorkLogID(workLogID string) (workrecordmodel.WorkLog, error) {
	return r.store.FindByWorkLogID(workLogID)
}

func (r *WorkLogRepository) FindByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.WorkLog, error) {
	return r.store.FindWorkLogsByUserIDAndDate(userID, workDate)
}

func (r *WorkLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.WorkLog, error) {
	return r.store.FindWorkLogsByUserIDAndWeek(userID, weekStartDate)
}

func (r *WorkLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.WorkLog, error) {
	return r.store.FindWorkLogsByUserIDAndMonth(userID, yearMonth)
}

type BreakLogRepository struct {
	store *memory.Store
}

func NewBreakLogRepository(store *memory.Store) *BreakLogRepository {
	return &BreakLogRepository{store: store}
}

func (r *BreakLogRepository) Create(breakLog workrecordmodel.BreakLog) error {
	return r.store.CreateBreakLog(breakLog)
}

func (r *BreakLogRepository) Update(breakLog workrecordmodel.BreakLog) error {
	return r.store.UpdateBreakLog(breakLog)
}

func (r *BreakLogRepository) Delete(breakLogID string) error {
	return r.store.DeleteBreakLog(breakLogID)
}

func (r *BreakLogRepository) FindByBreakLogID(breakLogID string) (workrecordmodel.BreakLog, error) {
	return r.store.FindByBreakLogID(breakLogID)
}

func (r *BreakLogRepository) FindByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.BreakLog, error) {
	return r.store.FindBreakLogsByUserIDAndDate(userID, workDate)
}

func (r *BreakLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.BreakLog, error) {
	return r.store.FindBreakLogsByUserIDAndWeek(userID, weekStartDate)
}

func (r *BreakLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.BreakLog, error) {
	return r.store.FindBreakLogsByUserIDAndMonth(userID, yearMonth)
}
