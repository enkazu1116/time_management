package repositories

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
)

type DailyAchievementRepository interface {
	Save(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error
	Update(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error
	Delete(dailyAchievementID string) error
	FindByUserIDAndDate(userID string, workDate time.Time) (analysismodel.DailySummarySet, error)
}
