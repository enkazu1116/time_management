package services

import (
	"time"

	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func BuildDailyBreakSummary(dailyAchievementID string, breakLogs []workrecordmodel.BreakLog) analysismodel.DailyBreakSummary {
	totalBreakDuration, averageBreakDuration, averageResumeDuration, longestResumeDuration := AnalyzeBreakLogs(breakLogs)
	return analysismodel.NewDailyBreakSummary(dailyAchievementID, len(breakLogs), totalBreakDuration, averageBreakDuration, averageResumeDuration, longestResumeDuration)
}

func BuildMonthlyBreakSummary(monthlyAchievementID string, breakLogs []workrecordmodel.BreakLog) analysismodel.MonthlyBreakSummary {
	totalBreakDuration, averageBreakDuration, averageResumeDuration, longestResumeDuration := AnalyzeBreakLogs(breakLogs)
	return analysismodel.NewMonthlyBreakSummary(monthlyAchievementID, len(breakLogs), totalBreakDuration, averageBreakDuration, averageResumeDuration, longestResumeDuration)
}

func AnalyzeBreakLogs(breakLogs []workrecordmodel.BreakLog) (time.Duration, time.Duration, time.Duration, time.Duration) {
	totalBreakDuration := time.Duration(0)
	totalResumeDuration := time.Duration(0)
	longestResumeDuration := time.Duration(0)
	resumeCount := 0
	for _, breakLog := range breakLogs {
		totalBreakDuration += breakLog.BreakDuration
		if breakLog.ResumeDuration != nil {
			totalResumeDuration += *breakLog.ResumeDuration
			resumeCount++
			if *breakLog.ResumeDuration > longestResumeDuration {
				longestResumeDuration = *breakLog.ResumeDuration
			}
		}
	}
	return totalBreakDuration, AverageDuration(totalBreakDuration, len(breakLogs)), AverageDuration(totalResumeDuration, resumeCount), longestResumeDuration
}
