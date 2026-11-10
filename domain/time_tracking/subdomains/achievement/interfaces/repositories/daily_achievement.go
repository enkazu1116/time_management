package repositories

import (
	"errors"
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
)

var ErrNotFound = errors.New("daily achievement not found")

type DailyAchievementRepository interface {
	Save(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error
	Update(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error
	Delete(dailyAchievementID string) error
	FindByUserIDAndDate(userID string, workDate time.Time) (analysismodel.DailySummarySet, error)
}
