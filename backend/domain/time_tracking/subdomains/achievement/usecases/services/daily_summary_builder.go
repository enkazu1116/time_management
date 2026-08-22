package services

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	analysisservices "time_management/domain/time_tracking/subdomains/analysis/usecases/services"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func BuildDailySummarySet(dailyAchievementID string, userID string, workDate time.Time, workLogs []workrecordmodel.WorkLog, breakLogs []workrecordmodel.BreakLog, updatedAt time.Time) analysismodel.DailySummarySet {
	totalWorkDuration := time.Duration(0)
	workTypeSummaries := map[string]analysismodel.DailyWorkTypeSummary{}
	for _, workLog := range workLogs {
		totalWorkDuration += workLog.WorkDuration
		summary := workTypeSummaries[workLog.WorkType.WorkTypeID]
		summary.DailyAchievementID = dailyAchievementID
		summary.WorkTypeID = workLog.WorkType.WorkTypeID
		summary.WorkCount++
		summary.WorkDuration += workLog.WorkDuration
		workTypeSummaries[workLog.WorkType.WorkTypeID] = summary
	}
	records := make([]analysismodel.DailyWorkTypeSummary, 0, len(workTypeSummaries))
	for _, summary := range workTypeSummaries {
		summary.AverageWorkDuration = analysisservices.AverageDuration(summary.WorkDuration, summary.WorkCount)
		records = append(records, summary)
	}
	breakSummary := analysisservices.BuildDailyBreakSummary(dailyAchievementID, breakLogs)
	return analysismodel.DailySummarySet{
		Achievement:       achievementmodel.NewDailyAchievement(dailyAchievementID, userID, workDate, totalWorkDuration, breakSummary.BreakDuration, len(workLogs), len(breakLogs), updatedAt),
		WorkTypeSummaries: records,
		BreakSummary:      breakSummary,
	}
}
