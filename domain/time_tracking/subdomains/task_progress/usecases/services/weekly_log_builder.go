package services

import (
	"time"

	sharedmodel "time_management/domain/shared/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	analysisservices "time_management/domain/time_tracking/subdomains/analysis/usecases/services"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

func BuildWeeklyLogSet(weeklyLogIDPrefix string, weeklyTasks []taskprogressmodel.WeeklyTask, dailyLogs []taskprogressmodel.DailyLog, workLogs []workrecordmodel.WorkLog, updatedAt time.Time) analysismodel.WeeklyLogSet {
	weeklyLogs := make([]taskprogressmodel.WeeklyLog, 0, len(weeklyTasks))
	for _, weeklyTask := range weeklyTasks {
		actualDuration := actualDurationByWorkType(weeklyTask.WorkTypeID, workLogs)
		averageProgressRate, progressEvaluation := analysisservices.EvaluateProgress(weeklyTask.WeeklyTaskID, dailyLogs)
		weeklyLogID := sharedmodel.NewDeterministicUUID(weeklyLogIDPrefix + "-" + weeklyTask.WeeklyTaskID).String()
		weeklyLogs = append(weeklyLogs, taskprogressmodel.NewWeeklyLog(
			weeklyLogID,
			weeklyTask.WeeklyTaskID,
			weeklyTask.UserID,
			weeklyTask.WorkTypeID,
			analysisservices.NormalizeWeekStart(weeklyTask.WeekStartDate),
			weeklyTask.TargetDuration,
			actualDuration,
			analysisservices.AchievementRate(actualDuration, weeklyTask.TargetDuration),
			weeklyTask.GoalDescription,
			progressEvaluation,
			averageProgressRate,
			updatedAt,
		))
	}
	return analysismodel.WeeklyLogSet{WeeklyLogs: weeklyLogs}
}

func actualDurationByWorkType(workTypeID string, workLogs []workrecordmodel.WorkLog) time.Duration {
	total := time.Duration(0)
	for _, workLog := range workLogs {
		if workLog.WorkType.WorkTypeID == workTypeID {
			total += workLog.WorkDuration
		}
	}
	return total
}
