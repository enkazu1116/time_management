package repositories

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
)

type MonthlyAchievementRepository interface {
	Save(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error
	Update(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error
	Delete(monthlyAchievementID string) error
	FindByUserIDAndMonth(userID string, yearMonth time.Time) (analysismodel.MonthlySummarySet, error)
}
