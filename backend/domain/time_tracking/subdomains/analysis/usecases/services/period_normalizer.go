package services

import "time"

func NormalizeYearMonth(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), 1, 0, 0, 0, 0, value.Location())
}

func NormalizeWeekStart(value time.Time) time.Time {
	weekdayOffset := (int(value.Weekday()) + 6) % 7
	weekStart := value.AddDate(0, 0, -weekdayOffset)
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, value.Location())
}

func AverageDuration(total time.Duration, count int) time.Duration {
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
