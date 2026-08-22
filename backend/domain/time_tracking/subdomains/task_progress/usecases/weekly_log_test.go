package usecases

import (
	"strings"
	"testing"
	"time"

	taskprogressrepositories "time_management/domain/time_tracking/subdomains/task_progress/interfaces/repositories"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	taskprogressports "time_management/domain/time_tracking/subdomains/task_progress/ports"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func TestCreateWeeklyLog(t *testing.T) {
	weekStart := at(2026, 8, 17, 0, 0)
	study := workrecordmodel.NewWorkType("00000000-0000-4000-8000-000000000401", "study")
	weeklyTasks := []taskprogressmodel.WeeklyTask{
		taskprogressmodel.NewWeeklyTask("00000000-0000-4000-8000-000000000701", "00000000-0000-4000-8000-000000000101", "00000000-0000-4000-8000-000000000401", weekStart, 10*time.Hour, "chapter 1 complete"),
	}
	dailyLogs := []taskprogressmodel.DailyLog{
		taskprogressmodel.NewDailyLog("00000000-0000-4000-8000-000000000901", "00000000-0000-4000-8000-000000000101", weekStart, 2*time.Hour, []taskprogressmodel.DailyTask{
			taskprogressmodel.NewDailyTask("00000000-0000-4000-8000-000000000801", "00000000-0000-4000-8000-000000000701", "00000000-0000-4000-8000-000000000101", weekStart, "read half", 0.5),
		}),
		taskprogressmodel.NewDailyLog("00000000-0000-4000-8000-000000000902", "00000000-0000-4000-8000-000000000101", at(2026, 8, 18, 0, 0), 3*time.Hour, []taskprogressmodel.DailyTask{
			taskprogressmodel.NewDailyTask("00000000-0000-4000-8000-000000000802", "00000000-0000-4000-8000-000000000701", "00000000-0000-4000-8000-000000000101", at(2026, 8, 18, 0, 0), "finished chapter", 1.0),
		}),
	}
	workLogs := []workrecordmodel.WorkLog{
		workrecordmodel.NewWorkLog("00000000-0000-4000-8000-000000000501", "00000000-0000-4000-8000-000000000101", study, weekStart, at(2026, 8, 17, 8, 0), at(2026, 8, 17, 10, 0), 2*time.Hour),
		workrecordmodel.NewWorkLog("00000000-0000-4000-8000-000000000502", "00000000-0000-4000-8000-000000000101", study, at(2026, 8, 18, 0, 0), at(2026, 8, 18, 8, 0), at(2026, 8, 18, 11, 0), 3*time.Hour),
	}
	weeklyLogRepo := &stubWeeklyLogRepository{}
	uc := NewWeeklyLogUsecase(&stubWorkLogRepository{weeklyLogs: workLogs}, &stubWeeklyTaskRepository{weeklyTasks: weeklyTasks}, &stubDailyLogRepository{weeklyLogs: dailyLogs}, weeklyLogRepo)

	weeklyLogSet, err := uc.CreateWeeklyLog(taskprogressports.CreateWeeklyLogInput{
		WeeklyLogIDPrefix: "weekly-log",
		UserID:            "00000000-0000-4000-8000-000000000101",
		WeekStartDate:     weekStart,
		UpdatedAt:         at(2026, 8, 23, 23, 0),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(weeklyLogSet.WeeklyLogs) != 1 {
		t.Fatalf("expected 1 weekly log, got %d", len(weeklyLogSet.WeeklyLogs))
	}
	weeklyLog := weeklyLogSet.WeeklyLogs[0]
	if weeklyLog.ActualDuration != 5*time.Hour {
		t.Fatalf("expected actual duration 5h, got %s", weeklyLog.ActualDuration)
	}
	if weeklyLog.AchievementRate != 0.5 {
		t.Fatalf("expected achievement rate 0.5, got %f", weeklyLog.AchievementRate)
	}
	if weeklyLog.AverageProgressRate != 0.75 {
		t.Fatalf("expected average progress rate 0.75, got %f", weeklyLog.AverageProgressRate)
	}
	if !strings.Contains(weeklyLog.ProgressEvaluation, "finished chapter") {
		t.Fatalf("expected progress evaluation to include daily qualitative progress")
	}
	if !weeklyLogRepo.saved {
		t.Fatalf("expected weekly logs to be saved")
	}
}

type stubWorkLogRepository struct {
	weeklyLogs []workrecordmodel.WorkLog
}

func (r *stubWorkLogRepository) Create(workLog workrecordmodel.WorkLog) error { return nil }
func (r *stubWorkLogRepository) Update(workLog workrecordmodel.WorkLog) error { return nil }
func (r *stubWorkLogRepository) Delete(workLogID string) error                { return nil }
func (r *stubWorkLogRepository) FindByWorkLogID(workLogID string) (workrecordmodel.WorkLog, error) {
	return workrecordmodel.WorkLog{}, nil
}
func (r *stubWorkLogRepository) FindByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.WorkLog, error) {
	return nil, nil
}
func (r *stubWorkLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.WorkLog, error) {
	return r.weeklyLogs, nil
}
func (r *stubWorkLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.WorkLog, error) {
	return nil, nil
}

type stubWeeklyTaskRepository struct {
	weeklyTasks []taskprogressmodel.WeeklyTask
}

func (r *stubWeeklyTaskRepository) Create(weeklyTask taskprogressmodel.WeeklyTask) error { return nil }
func (r *stubWeeklyTaskRepository) Update(weeklyTask taskprogressmodel.WeeklyTask) error { return nil }
func (r *stubWeeklyTaskRepository) Delete(weeklyTaskID string) error                     { return nil }
func (r *stubWeeklyTaskRepository) FindByWeeklyTaskID(weeklyTaskID string) (taskprogressmodel.WeeklyTask, error) {
	return taskprogressmodel.WeeklyTask{}, nil
}
func (r *stubWeeklyTaskRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.WeeklyTask, error) {
	return r.weeklyTasks, nil
}

type stubDailyLogRepository struct {
	weeklyLogs []taskprogressmodel.DailyLog
}

func (r *stubDailyLogRepository) Create(dailyLog taskprogressmodel.DailyLog) error { return nil }
func (r *stubDailyLogRepository) Update(dailyLog taskprogressmodel.DailyLog) error { return nil }
func (r *stubDailyLogRepository) Delete(dailyLogID string) error                   { return nil }
func (r *stubDailyLogRepository) FindByDailyLogID(dailyLogID string) (taskprogressmodel.DailyLog, error) {
	return taskprogressmodel.DailyLog{}, nil
}
func (r *stubDailyLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.DailyLog, error) {
	return r.weeklyLogs, nil
}

type stubWeeklyLogRepository struct {
	saved       bool
	monthlyLogs []taskprogressmodel.WeeklyLog
}

func (r *stubWeeklyLogRepository) Create(weeklyLog taskprogressmodel.WeeklyLog) error { return nil }
func (r *stubWeeklyLogRepository) Update(weeklyLog taskprogressmodel.WeeklyLog) error { return nil }
func (r *stubWeeklyLogRepository) Delete(weeklyLogID string) error                    { return nil }
func (r *stubWeeklyLogRepository) SaveAll(weeklyLogs []taskprogressmodel.WeeklyLog) error {
	r.saved = true
	return nil
}
func (r *stubWeeklyLogRepository) FindByWeeklyLogID(weeklyLogID string) (taskprogressmodel.WeeklyLog, error) {
	return taskprogressmodel.WeeklyLog{}, nil
}
func (r *stubWeeklyLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]taskprogressmodel.WeeklyLog, error) {
	return r.monthlyLogs, nil
}

var _ taskprogressrepositories.WeeklyLogRepository = (*stubWeeklyLogRepository)(nil)

func at(year int, month time.Month, day int, hour int, minute int) time.Time {
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}
