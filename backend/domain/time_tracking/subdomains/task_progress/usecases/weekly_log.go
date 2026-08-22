package usecases

import (
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressrepositories "time_management/domain/time_tracking/subdomains/task_progress/interfaces/repositories"
	taskprogressports "time_management/domain/time_tracking/subdomains/task_progress/ports"
	taskprogressservices "time_management/domain/time_tracking/subdomains/task_progress/usecases/services"
	workrecordrepositories "time_management/domain/time_tracking/subdomains/work_record/interfaces/repositories"
)

type WeeklyLogUsecase interface {
	CreateWeeklyLog(input taskprogressports.CreateWeeklyLogInput) (analysismodel.WeeklyLogSet, error)
}

type weeklyLogUsecase struct {
	workLogRepository    workrecordrepositories.WorkLogRepository
	weeklyTaskRepository taskprogressrepositories.WeeklyTaskRepository
	dailyLogRepository   taskprogressrepositories.DailyLogRepository
	weeklyLogRepository  taskprogressrepositories.WeeklyLogRepository
}

func NewWeeklyLogUsecase(
	workLogRepository workrecordrepositories.WorkLogRepository,
	weeklyTaskRepository taskprogressrepositories.WeeklyTaskRepository,
	dailyLogRepository taskprogressrepositories.DailyLogRepository,
	weeklyLogRepository taskprogressrepositories.WeeklyLogRepository,
) WeeklyLogUsecase {
	return &weeklyLogUsecase{
		workLogRepository:    workLogRepository,
		weeklyTaskRepository: weeklyTaskRepository,
		dailyLogRepository:   dailyLogRepository,
		weeklyLogRepository:  weeklyLogRepository,
	}
}

func (u *weeklyLogUsecase) CreateWeeklyLog(input taskprogressports.CreateWeeklyLogInput) (analysismodel.WeeklyLogSet, error) {
	if err := input.Validate(); err != nil {
		return analysismodel.WeeklyLogSet{}, err
	}
	weeklyTasks, err := u.weeklyTaskRepository.FindByUserIDAndWeek(input.UserID, input.WeekStartDate)
	if err != nil {
		return analysismodel.WeeklyLogSet{}, err
	}
	dailyLogs, err := u.dailyLogRepository.FindByUserIDAndWeek(input.UserID, input.WeekStartDate)
	if err != nil {
		return analysismodel.WeeklyLogSet{}, err
	}
	workLogs, err := u.workLogRepository.FindByUserIDAndWeek(input.UserID, input.WeekStartDate)
	if err != nil {
		return analysismodel.WeeklyLogSet{}, err
	}
	weeklyLogSet := taskprogressservices.BuildWeeklyLogSet(input.WeeklyLogIDPrefix, weeklyTasks, dailyLogs, workLogs, input.UpdatedAt)
	for _, weeklyLog := range weeklyLogSet.WeeklyLogs {
		if err := weeklyLog.Validate(); err != nil {
			return analysismodel.WeeklyLogSet{}, err
		}
	}
	if err := u.weeklyLogRepository.SaveAll(weeklyLogSet.WeeklyLogs); err != nil {
		return analysismodel.WeeklyLogSet{}, err
	}
	return weeklyLogSet, nil
}
