package services

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	analysisservices "time_management/domain/time_tracking/subdomains/analysis/usecases/services"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func BuildMonthlySummarySetFromWeeklyLogs(monthlyAchievementID string, userID string, yearMonth time.Time, weeklyLogs []taskprogressmodel.WeeklyLog, updatedAt time.Time) analysismodel.MonthlySummarySet {
	month := analysisservices.NormalizeYearMonth(yearMonth)
	filteredWeeklyLogs := filterWeeklyLogsByMonth(weeklyLogs, month)
	totalWorkDuration := time.Duration(0)
	workTypeSummaries := map[string]analysismodel.MonthlyWorkTypeSummary{}
	workDaysByWeek := map[string]struct{}{}
	for _, weeklyLog := range filteredWeeklyLogs {
		totalWorkDuration += weeklyLog.ActualDuration
		workDaysByWeek[weeklyLog.WeekStartDate.Format("2006-01-02")] = struct{}{}
		summary := workTypeSummaries[weeklyLog.WorkTypeID]
		summary.MonthlyAchievementID = monthlyAchievementID
		summary.WorkTypeID = weeklyLog.WorkTypeID
		summary.WorkCount++
		summary.WorkDuration += weeklyLog.ActualDuration
		workTypeSummaries[weeklyLog.WorkTypeID] = summary
	}
	records := make([]analysismodel.MonthlyWorkTypeSummary, 0, len(workTypeSummaries))
	for _, summary := range workTypeSummaries {
		summary.AverageWorkDuration = analysisservices.AverageDuration(summary.WorkDuration, summary.WorkCount)
		records = append(records, summary)
	}
	return analysismodel.MonthlySummarySet{
		Achievement:       achievementmodel.NewMonthlyAchievement(monthlyAchievementID, userID, month, totalWorkDuration, 0, len(workDaysByWeek), len(filteredWeeklyLogs), 0, updatedAt),
		WorkTypeSummaries: records,
		BreakSummary:      analysismodel.NewMonthlyBreakSummary(monthlyAchievementID, 0, 0, 0, 0, 0),
	}
}

func BuildMonthlySummarySetFromWorkLogs(monthlyAchievementID string, userID string, yearMonth time.Time, workLogs []workrecordmodel.WorkLog, breakLogs []workrecordmodel.BreakLog, updatedAt time.Time) analysismodel.MonthlySummarySet {
	totalWorkDuration := time.Duration(0)
	workDaysByDate := map[string]struct{}{}
	workTypeSummaries := map[string]analysismodel.MonthlyWorkTypeSummary{}
	for _, workLog := range workLogs {
		totalWorkDuration += workLog.WorkDuration
		workDaysByDate[workLog.WorkDate.Format("2006-01-02")] = struct{}{}
		summary := workTypeSummaries[workLog.WorkType.WorkTypeID]
		summary.MonthlyAchievementID = monthlyAchievementID
		summary.WorkTypeID = workLog.WorkType.WorkTypeID
		summary.WorkCount++
		summary.WorkDuration += workLog.WorkDuration
		workTypeSummaries[workLog.WorkType.WorkTypeID] = summary
	}
	records := make([]analysismodel.MonthlyWorkTypeSummary, 0, len(workTypeSummaries))
	for _, summary := range workTypeSummaries {
		summary.AverageWorkDuration = analysisservices.AverageDuration(summary.WorkDuration, summary.WorkCount)
		records = append(records, summary)
	}
	breakSummary := analysisservices.BuildMonthlyBreakSummary(monthlyAchievementID, breakLogs)
	return analysismodel.MonthlySummarySet{
		Achievement:       achievementmodel.NewMonthlyAchievement(monthlyAchievementID, userID, analysisservices.NormalizeYearMonth(yearMonth), totalWorkDuration, breakSummary.BreakDuration, len(workDaysByDate), len(workLogs), len(breakLogs), updatedAt),
		WorkTypeSummaries: records,
		BreakSummary:      breakSummary,
	}
}

func filterWeeklyLogsByMonth(weeklyLogs []taskprogressmodel.WeeklyLog, yearMonth time.Time) []taskprogressmodel.WeeklyLog {
	filtered := make([]taskprogressmodel.WeeklyLog, 0, len(weeklyLogs))
	for _, weeklyLog := range weeklyLogs {
		if analysisservices.NormalizeYearMonth(weeklyLog.WeekStartDate).Equal(yearMonth) {
			filtered = append(filtered, weeklyLog)
		}
	}
	return filtered
}
