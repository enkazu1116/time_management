package services

import (
	"strings"
	"time"

	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
)

func NormalizeProgressRate(progressRate float64) float64 {
	if progressRate < 0 {
		return 0
	}
	if progressRate > 1 {
		return 1
	}
	return progressRate
}

func EvaluateProgress(weeklyTaskID string, dailyLogs []taskprogressmodel.DailyLog) (float64, string) {
	totalProgressRate := 0.0
	progressCount := 0
	progressDescriptions := []string{}
	for _, dailyLog := range dailyLogs {
		for _, dailyTask := range dailyLog.DailyTasks {
			if dailyTask.WeeklyTaskID != weeklyTaskID {
				continue
			}
			totalProgressRate += dailyTask.ProgressRate
			progressCount++
			if dailyTask.ProgressDescription != "" {
				progressDescriptions = append(progressDescriptions, dailyTask.ProgressDescription)
			}
		}
	}
	if progressCount == 0 {
		return 0, ""
	}
	return totalProgressRate / float64(progressCount), strings.Join(progressDescriptions, "\n")
}

func AchievementRate(actualDuration time.Duration, targetDuration time.Duration) float64 {
	if targetDuration == 0 {
		return 0
	}
	return float64(actualDuration) / float64(targetDuration)
}
