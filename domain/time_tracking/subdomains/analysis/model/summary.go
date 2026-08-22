package model

import (
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
)

type DailyWorkTypeSummary struct {
	DailyAchievementID  string
	WorkTypeID          string
	WorkCount           int
	WorkDuration        time.Duration
	AverageWorkDuration time.Duration
}

type DailyBreakSummary struct {
	DailyAchievementID    string
	BreakCount            int
	BreakDuration         time.Duration
	AverageBreakDuration  time.Duration
	AverageResumeDuration time.Duration
	LongestResumeDuration time.Duration
}

type MonthlyWorkTypeSummary struct {
	MonthlyAchievementID string
	WorkTypeID           string
	WorkCount            int
	WorkDuration         time.Duration
	AverageWorkDuration  time.Duration
}

type MonthlyBreakSummary struct {
	MonthlyAchievementID  string
	BreakCount            int
	BreakDuration         time.Duration
	AverageBreakDuration  time.Duration
	AverageResumeDuration time.Duration
	LongestResumeDuration time.Duration
}

type DailySummarySet struct {
	Achievement       achievementmodel.DailyAchievement
	WorkTypeSummaries []DailyWorkTypeSummary
	BreakSummary      DailyBreakSummary
}

type WeeklyLogSet struct {
	WeeklyLogs []taskprogressmodel.WeeklyLog
}

type MonthlySummarySet struct {
	Achievement       achievementmodel.MonthlyAchievement
	WorkTypeSummaries []MonthlyWorkTypeSummary
	BreakSummary      MonthlyBreakSummary
}

func NewDailyWorkTypeSummary(dailyAchievementID string, workTypeID string, workCount int, workDuration time.Duration, averageWorkDuration time.Duration) DailyWorkTypeSummary {
	return DailyWorkTypeSummary{DailyAchievementID: dailyAchievementID, WorkTypeID: workTypeID, WorkCount: workCount, WorkDuration: workDuration, AverageWorkDuration: averageWorkDuration}
}

func NewDailyBreakSummary(dailyAchievementID string, breakCount int, breakDuration time.Duration, averageBreakDuration time.Duration, averageResumeDuration time.Duration, longestResumeDuration time.Duration) DailyBreakSummary {
	return DailyBreakSummary{DailyAchievementID: dailyAchievementID, BreakCount: breakCount, BreakDuration: breakDuration, AverageBreakDuration: averageBreakDuration, AverageResumeDuration: averageResumeDuration, LongestResumeDuration: longestResumeDuration}
}

func NewMonthlyWorkTypeSummary(monthlyAchievementID string, workTypeID string, workCount int, workDuration time.Duration, averageWorkDuration time.Duration) MonthlyWorkTypeSummary {
	return MonthlyWorkTypeSummary{MonthlyAchievementID: monthlyAchievementID, WorkTypeID: workTypeID, WorkCount: workCount, WorkDuration: workDuration, AverageWorkDuration: averageWorkDuration}
}

func NewMonthlyBreakSummary(monthlyAchievementID string, breakCount int, breakDuration time.Duration, averageBreakDuration time.Duration, averageResumeDuration time.Duration, longestResumeDuration time.Duration) MonthlyBreakSummary {
	return MonthlyBreakSummary{MonthlyAchievementID: monthlyAchievementID, BreakCount: breakCount, BreakDuration: breakDuration, AverageBreakDuration: averageBreakDuration, AverageResumeDuration: averageResumeDuration, LongestResumeDuration: longestResumeDuration}
}
