package model

import "time"

type DailyAchievement struct {
	DailyAchievementID string
	UserID             string
	WorkDate           time.Time
	TotalWorkDuration  time.Duration
	TotalBreakDuration time.Duration
	WorkCount          int
	BreakCount         int
	UpdatedAt          time.Time
}

type MonthlyAchievement struct {
	MonthlyAchievementID string
	UserID               string
	YearMonth            time.Time
	TotalWorkDuration    time.Duration
	TotalBreakDuration   time.Duration
	WorkDays             int
	WorkCount            int
	BreakCount           int
	UpdatedAt            time.Time
}

func NewDailyAchievement(dailyAchievementID string, userID string, workDate time.Time, totalWorkDuration time.Duration, totalBreakDuration time.Duration, workCount int, breakCount int, updatedAt time.Time) DailyAchievement {
	return DailyAchievement{DailyAchievementID: dailyAchievementID, UserID: userID, WorkDate: workDate, TotalWorkDuration: totalWorkDuration, TotalBreakDuration: totalBreakDuration, WorkCount: workCount, BreakCount: breakCount, UpdatedAt: updatedAt}
}

func NewMonthlyAchievement(monthlyAchievementID string, userID string, yearMonth time.Time, totalWorkDuration time.Duration, totalBreakDuration time.Duration, workDays int, workCount int, breakCount int, updatedAt time.Time) MonthlyAchievement {
	return MonthlyAchievement{MonthlyAchievementID: monthlyAchievementID, UserID: userID, YearMonth: yearMonth, TotalWorkDuration: totalWorkDuration, TotalBreakDuration: totalBreakDuration, WorkDays: workDays, WorkCount: workCount, BreakCount: breakCount, UpdatedAt: updatedAt}
}
