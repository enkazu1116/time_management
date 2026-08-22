package usecases

import (
	"testing"
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	achievementports "time_management/domain/time_tracking/subdomains/achievement/ports"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func TestCreateMonthlySummaryFromWeeklyLogs(t *testing.T) {
	weeklyLogs := []taskprogressmodel.WeeklyLog{
		taskprogressmodel.NewWeeklyLog("00000000-0000-4000-8000-000000000d01", "00000000-0000-4000-8000-000000000a01", "00000000-0000-4000-8000-000000000101", "00000000-0000-4000-8000-000000000401", at(2026, 8, 3, 0, 0), 10*time.Hour, 8*time.Hour, 0.8, "study", "good", 0.8, at(2026, 8, 9, 23, 0)),
		taskprogressmodel.NewWeeklyLog("00000000-0000-4000-8000-000000000d02", "00000000-0000-4000-8000-000000000a02", "00000000-0000-4000-8000-000000000101", "00000000-0000-4000-8000-000000000402", at(2026, 8, 10, 0, 0), 3*time.Hour, 2*time.Hour, 0.66, "workout", "done", 1.0, at(2026, 8, 16, 23, 0)),
		taskprogressmodel.NewWeeklyLog("00000000-0000-4000-8000-000000000d03", "00000000-0000-4000-8000-000000000a03", "00000000-0000-4000-8000-000000000101", "00000000-0000-4000-8000-000000000401", at(2026, 9, 7, 0, 0), 10*time.Hour, 6*time.Hour, 0.6, "study", "next month", 0.6, at(2026, 9, 13, 23, 0)),
	}
	monthlyAchievementRepo := &stubMonthlyAchievementRepository{}
	uc := NewSummaryUsecase(&stubWorkLogRepository{}, &stubBreakLogRepository{}, &stubDailyAchievementRepository{}, monthlyAchievementRepo, &stubWeeklyLogRepository{monthlyLogs: weeklyLogs})

	monthly, err := uc.CreateMonthlySummary(achievementports.CreateMonthlySummaryInput{
		MonthlyAchievementID: "00000000-0000-4000-8000-000000000c01",
		UserID:               "00000000-0000-4000-8000-000000000101",
		YearMonth:            at(2026, 8, 25, 12, 0),
		UpdatedAt:            at(2026, 8, 31, 23, 0),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if monthly.Achievement.TotalWorkDuration != 10*time.Hour {
		t.Fatalf("expected total work duration 10h, got %s", monthly.Achievement.TotalWorkDuration)
	}
	if monthly.Achievement.WorkCount != 2 {
		t.Fatalf("expected only August weekly logs, got %d", monthly.Achievement.WorkCount)
	}
	if !monthlyAchievementRepo.saved {
		t.Fatalf("expected monthly summary to be saved")
	}
}

type stubWorkLogRepository struct{}

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
	return nil, nil
}
func (r *stubWorkLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.WorkLog, error) {
	return nil, nil
}

type stubBreakLogRepository struct{}

func (r *stubBreakLogRepository) Create(breakLog workrecordmodel.BreakLog) error { return nil }
func (r *stubBreakLogRepository) Update(breakLog workrecordmodel.BreakLog) error { return nil }
func (r *stubBreakLogRepository) Delete(breakLogID string) error                 { return nil }
func (r *stubBreakLogRepository) FindByBreakLogID(breakLogID string) (workrecordmodel.BreakLog, error) {
	return workrecordmodel.BreakLog{}, nil
}
func (r *stubBreakLogRepository) FindByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.BreakLog, error) {
	return nil, nil
}
func (r *stubBreakLogRepository) FindByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.BreakLog, error) {
	return nil, nil
}
func (r *stubBreakLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.BreakLog, error) {
	return nil, nil
}

type stubDailyAchievementRepository struct{}

func (r *stubDailyAchievementRepository) Save(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	return nil
}
func (r *stubDailyAchievementRepository) Update(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	return nil
}
func (r *stubDailyAchievementRepository) Delete(dailyAchievementID string) error { return nil }
func (r *stubDailyAchievementRepository) FindByUserIDAndDate(userID string, workDate time.Time) (analysismodel.DailySummarySet, error) {
	return analysismodel.DailySummarySet{}, nil
}

type stubMonthlyAchievementRepository struct {
	saved bool
}

func (r *stubMonthlyAchievementRepository) Save(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	r.saved = true
	return nil
}
func (r *stubMonthlyAchievementRepository) Update(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	return nil
}
func (r *stubMonthlyAchievementRepository) Delete(monthlyAchievementID string) error { return nil }
func (r *stubMonthlyAchievementRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) (analysismodel.MonthlySummarySet, error) {
	return analysismodel.MonthlySummarySet{}, nil
}

type stubWeeklyLogRepository struct {
	monthlyLogs []taskprogressmodel.WeeklyLog
}

func (r *stubWeeklyLogRepository) Create(weeklyLog taskprogressmodel.WeeklyLog) error { return nil }
func (r *stubWeeklyLogRepository) Update(weeklyLog taskprogressmodel.WeeklyLog) error { return nil }
func (r *stubWeeklyLogRepository) Delete(weeklyLogID string) error                    { return nil }
func (r *stubWeeklyLogRepository) SaveAll(weeklyLogs []taskprogressmodel.WeeklyLog) error {
	return nil
}
func (r *stubWeeklyLogRepository) FindByWeeklyLogID(weeklyLogID string) (taskprogressmodel.WeeklyLog, error) {
	return taskprogressmodel.WeeklyLog{}, nil
}
func (r *stubWeeklyLogRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) ([]taskprogressmodel.WeeklyLog, error) {
	return r.monthlyLogs, nil
}

func at(year int, month time.Month, day int, hour int, minute int) time.Time {
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}
