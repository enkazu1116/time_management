package achievement

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	"time_management/infrastructure/persistence/memory"
)

type DailyAchievementRepository struct {
	store *memory.Store
}

func NewDailyAchievementRepository(store *memory.Store) *DailyAchievementRepository {
	return &DailyAchievementRepository{store: store}
}

func (r *DailyAchievementRepository) Save(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	return r.store.SaveDaily(achievement, workTypeSummaries, breakSummary)
}

func (r *DailyAchievementRepository) Update(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	return r.store.UpdateDaily(achievement, workTypeSummaries, breakSummary)
}

func (r *DailyAchievementRepository) Delete(dailyAchievementID string) error {
	return r.store.DeleteDaily(dailyAchievementID)
}

func (r *DailyAchievementRepository) FindByUserIDAndDate(userID string, workDate time.Time) (analysismodel.DailySummarySet, error) {
	return r.store.FindDailyByUserIDAndDate(userID, workDate)
}

type MonthlyAchievementRepository struct {
	store *memory.Store
}

func NewMonthlyAchievementRepository(store *memory.Store) *MonthlyAchievementRepository {
	return &MonthlyAchievementRepository{store: store}
}

func (r *MonthlyAchievementRepository) Save(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	return r.store.SaveMonthly(achievement, workTypeSummaries, breakSummary)
}

func (r *MonthlyAchievementRepository) Update(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	return r.store.UpdateMonthly(achievement, workTypeSummaries, breakSummary)
}

func (r *MonthlyAchievementRepository) Delete(monthlyAchievementID string) error {
	return r.store.DeleteMonthly(monthlyAchievementID)
}

func (r *MonthlyAchievementRepository) FindByUserIDAndMonth(userID string, yearMonth time.Time) (analysismodel.MonthlySummarySet, error) {
	return r.store.FindMonthlyByUserIDAndMonth(userID, yearMonth)
}
