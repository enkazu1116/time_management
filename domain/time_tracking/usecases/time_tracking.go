package usecases

import (
	trackingports "time_management/domain/time_tracking/ports"
	achievementrepositories "time_management/domain/time_tracking/subdomains/achievement/interfaces/repositories"
	achievementports "time_management/domain/time_tracking/subdomains/achievement/ports"
	achievementusecases "time_management/domain/time_tracking/subdomains/achievement/usecases"
	taskprogressrepositories "time_management/domain/time_tracking/subdomains/task_progress/interfaces/repositories"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	taskprogressports "time_management/domain/time_tracking/subdomains/task_progress/ports"
	taskprogressusecases "time_management/domain/time_tracking/subdomains/task_progress/usecases"
	workrecordrepositories "time_management/domain/time_tracking/subdomains/work_record/interfaces/repositories"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

type Service struct {
	workLogs      workrecordrepositories.WorkLogRepository
	breakLogs     workrecordrepositories.BreakLogRepository
	weeklyTasks   taskprogressrepositories.WeeklyTaskRepository
	dailyTasks    taskprogressrepositories.DailyTaskRepository
	dailyLogs     taskprogressrepositories.DailyLogRepository
	weeklyLogs    taskprogressrepositories.WeeklyLogRepository
	summaries     achievementusecases.SummaryUsecase
	weeklySummary taskprogressusecases.WeeklyLogUsecase
}

func NewService(w workrecordrepositories.WorkLogRepository, b workrecordrepositories.BreakLogRepository, wt taskprogressrepositories.WeeklyTaskRepository, dt taskprogressrepositories.DailyTaskRepository, dl taskprogressrepositories.DailyLogRepository, wl taskprogressrepositories.WeeklyLogRepository, da achievementrepositories.DailyAchievementRepository, ma achievementrepositories.MonthlyAchievementRepository) *Service {
	return &Service{workLogs: w, breakLogs: b, weeklyTasks: wt, dailyTasks: dt, dailyLogs: dl, weeklyLogs: wl, summaries: achievementusecases.NewSummaryUsecase(w, b, da, ma, wl), weeklySummary: taskprogressusecases.NewWeeklyLogUsecase(w, wt, dl, wl)}
}

func (s *Service) CreateWorkLog(i trackingports.CreateWorkLogInput) (trackingports.WorkLogOutput, error) {
	v := workrecordmodel.NewWorkLog(i.WorkLogID, i.UserID, workrecordmodel.NewWorkType(i.WorkTypeID, i.WorkTypeName), i.WorkDate, i.StartedAt, i.EndedAt, i.WorkDuration)
	if err := v.Validate(); err != nil {
		return trackingports.WorkLogOutput{}, err
	}
	if err := s.workLogs.Create(v); err != nil {
		return trackingports.WorkLogOutput{}, err
	}
	return workLogOutput(v), nil
}
func (s *Service) CreateBreakLog(i trackingports.CreateBreakLogInput) (trackingports.BreakLogOutput, error) {
	v := workrecordmodel.NewBreakLog(i.BreakLogID, i.WorkLogID, i.UserID, i.StartedAt, i.EndedAt, i.BreakDuration, i.ResumedAt, i.ResumeDuration)
	if err := v.Validate(); err != nil {
		return trackingports.BreakLogOutput{}, err
	}
	if err := s.breakLogs.Create(v); err != nil {
		return trackingports.BreakLogOutput{}, err
	}
	return breakLogOutput(v), nil
}
func (s *Service) CreateWeeklyTask(i trackingports.CreateWeeklyTaskInput) (trackingports.WeeklyTaskOutput, error) {
	v := taskprogressmodel.NewWeeklyTask(i.WeeklyTaskID, i.UserID, i.WorkTypeID, i.WeekStartDate, i.TargetDuration, i.GoalDescription)
	if err := v.Validate(); err != nil {
		return trackingports.WeeklyTaskOutput{}, err
	}
	if err := s.weeklyTasks.Create(v); err != nil {
		return trackingports.WeeklyTaskOutput{}, err
	}
	return trackingports.WeeklyTaskOutput{v.WeeklyTaskID, v.UserID, v.WorkTypeID, v.WeekStartDate, v.TargetDuration, v.GoalDescription}, nil
}
func (s *Service) CreateDailyTask(i trackingports.CreateDailyTaskInput) (trackingports.DailyTaskOutput, error) {
	v := taskprogressmodel.NewDailyTask(i.DailyTaskID, i.WeeklyTaskID, i.UserID, i.WorkDate, i.ProgressDescription, i.ProgressRate)
	if err := v.Validate(); err != nil {
		return trackingports.DailyTaskOutput{}, err
	}
	if err := s.dailyTasks.Create(v); err != nil {
		return trackingports.DailyTaskOutput{}, err
	}
	return trackingports.DailyTaskOutput{v.DailyTaskID, v.WeeklyTaskID, v.UserID, v.WorkDate, v.ProgressDescription, v.ProgressRate}, nil
}
func (s *Service) CreateDailyLog(i trackingports.CreateDailyLogInput) (trackingports.DailyLogOutput, error) {
	tasks, err := s.dailyTasks.FindByUserIDAndDate(i.UserID, i.WorkDate)
	if err != nil {
		return trackingports.DailyLogOutput{}, err
	}
	v := taskprogressmodel.NewDailyLog(i.DailyLogID, i.UserID, i.WorkDate, i.TotalWorkDuration, tasks)
	if err := v.Validate(); err != nil {
		return trackingports.DailyLogOutput{}, err
	}
	if err := s.dailyLogs.Create(v); err != nil {
		return trackingports.DailyLogOutput{}, err
	}
	return trackingports.DailyLogOutput{v.DailyLogID, v.UserID, v.WorkDate, v.TotalWorkDuration}, nil
}
func (s *Service) CreateDailySummary(i trackingports.CreateDailySummaryInput) (trackingports.SummaryOutput, error) {
	v, err := s.summaries.CreateDailySummary(achievementports.CreateDailySummaryInput{DailyAchievementID: i.DailyAchievementID, UserID: i.UserID, WorkDate: i.WorkDate, UpdatedAt: i.UpdatedAt})
	return trackingports.SummaryOutput{Kind: "daily", AchievementID: v.Achievement.DailyAchievementID, UserID: v.Achievement.UserID, Period: v.Achievement.WorkDate.Format("2006-01-02"), TotalWorkDuration: v.Achievement.TotalWorkDuration, TotalBreakDuration: v.Achievement.TotalBreakDuration, WorkCount: v.Achievement.WorkCount, BreakCount: v.Achievement.BreakCount}, err
}
func (s *Service) CreateWeeklyLog(i trackingports.CreateWeeklyLogInput) (trackingports.WeeklyLogSetOutput, error) {
	v, err := s.weeklySummary.CreateWeeklyLog(taskprogressports.CreateWeeklyLogInput{WeeklyLogIDPrefix: i.WeeklyLogIDPrefix, UserID: i.UserID, WeekStartDate: i.WeekStartDate, UpdatedAt: i.UpdatedAt})
	if err != nil {
		return trackingports.WeeklyLogSetOutput{}, err
	}
	out := trackingports.WeeklyLogSetOutput{}
	for _, x := range v.WeeklyLogs {
		out.WeeklyLogs = append(out.WeeklyLogs, trackingports.WeeklyLogOutput{x.WeeklyLogID, x.WeeklyTaskID, x.UserID, x.WorkTypeID, x.WeekStartDate, x.TargetDuration, x.ActualDuration, x.AchievementRate, x.ProgressEvaluation})
	}
	return out, nil
}
func (s *Service) CreateMonthlySummary(i trackingports.CreateMonthlySummaryInput, fromWorkLogs bool) (trackingports.SummaryOutput, error) {
	v, err := s.summaries.CreateMonthlySummary(achievementports.CreateMonthlySummaryInput{MonthlyAchievementID: i.MonthlyAchievementID, UserID: i.UserID, YearMonth: i.YearMonth, UpdatedAt: i.UpdatedAt})
	if fromWorkLogs {
		v, err = s.summaries.CreateMonthlySummaryFromWorkLogs(achievementports.CreateMonthlySummaryInput{MonthlyAchievementID: i.MonthlyAchievementID, UserID: i.UserID, YearMonth: i.YearMonth, UpdatedAt: i.UpdatedAt})
	}
	return trackingports.SummaryOutput{Kind: "monthly", AchievementID: v.Achievement.MonthlyAchievementID, UserID: v.Achievement.UserID, Period: v.Achievement.YearMonth.Format("2006-01"), TotalWorkDuration: v.Achievement.TotalWorkDuration, TotalBreakDuration: v.Achievement.TotalBreakDuration, WorkDays: v.Achievement.WorkDays, WorkCount: v.Achievement.WorkCount, BreakCount: v.Achievement.BreakCount}, err
}

func workLogOutput(v workrecordmodel.WorkLog) trackingports.WorkLogOutput {
	return trackingports.WorkLogOutput{v.WorkLogID, v.UserID, v.WorkType.WorkTypeID, v.WorkType.WorkTypeName, v.WorkDate, v.StartedAt, v.EndedAt, v.WorkDuration}
}
func breakLogOutput(v workrecordmodel.BreakLog) trackingports.BreakLogOutput {
	return trackingports.BreakLogOutput{v.BreakLogID, v.WorkLogID, v.UserID, v.StartedAt, v.EndedAt, v.BreakDuration, v.ResumedAt, v.ResumeDuration}
}

var _ trackingports.TimeTrackingUsecase = (*Service)(nil)
