package usecases

import (
	achievementrepositories "time_management/domain/time_tracking/subdomains/achievement/interfaces/repositories"
	achievementports "time_management/domain/time_tracking/subdomains/achievement/ports"
	achievementservices "time_management/domain/time_tracking/subdomains/achievement/usecases/services"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressrepositories "time_management/domain/time_tracking/subdomains/task_progress/interfaces/repositories"
	workrecordrepositories "time_management/domain/time_tracking/subdomains/work_record/interfaces/repositories"
)

type SummaryUsecase interface {
	CreateDailySummary(input achievementports.CreateDailySummaryInput) (analysismodel.DailySummarySet, error)
	CreateMonthlySummary(input achievementports.CreateMonthlySummaryInput) (analysismodel.MonthlySummarySet, error)
	CreateMonthlySummaryFromWorkLogs(input achievementports.CreateMonthlySummaryInput) (analysismodel.MonthlySummarySet, error)
}

type summaryUsecase struct {
	workLogRepository      workrecordrepositories.WorkLogRepository
	breakLogRepository     workrecordrepositories.BreakLogRepository
	dailyAchievementRepo   achievementrepositories.DailyAchievementRepository
	monthlyAchievementRepo achievementrepositories.MonthlyAchievementRepository
	weeklyLogRepository    taskprogressrepositories.WeeklyLogRepository
}

func NewSummaryUsecase(
	workLogRepository workrecordrepositories.WorkLogRepository,
	breakLogRepository workrecordrepositories.BreakLogRepository,
	dailyAchievementRepo achievementrepositories.DailyAchievementRepository,
	monthlyAchievementRepo achievementrepositories.MonthlyAchievementRepository,
	weeklyLogRepository taskprogressrepositories.WeeklyLogRepository,
) SummaryUsecase {
	return &summaryUsecase{
		workLogRepository:      workLogRepository,
		breakLogRepository:     breakLogRepository,
		dailyAchievementRepo:   dailyAchievementRepo,
		monthlyAchievementRepo: monthlyAchievementRepo,
		weeklyLogRepository:    weeklyLogRepository,
	}
}

func (u *summaryUsecase) CreateDailySummary(input achievementports.CreateDailySummaryInput) (analysismodel.DailySummarySet, error) {
	if err := input.Validate(); err != nil {
		return analysismodel.DailySummarySet{}, err
	}
	workLogs, err := u.workLogRepository.FindByUserIDAndDate(input.UserID, input.WorkDate)
	if err != nil {
		return analysismodel.DailySummarySet{}, err
	}
	breakLogs, err := u.breakLogRepository.FindByUserIDAndDate(input.UserID, input.WorkDate)
	if err != nil {
		return analysismodel.DailySummarySet{}, err
	}
	summarySet := achievementservices.BuildDailySummarySet(input.DailyAchievementID, input.UserID, input.WorkDate, workLogs, breakLogs, input.UpdatedAt)
	if err := u.dailyAchievementRepo.Save(summarySet.Achievement, summarySet.WorkTypeSummaries, summarySet.BreakSummary); err != nil {
		return analysismodel.DailySummarySet{}, err
	}
	return summarySet, nil
}

func (u *summaryUsecase) CreateMonthlySummary(input achievementports.CreateMonthlySummaryInput) (analysismodel.MonthlySummarySet, error) {
	if err := input.Validate(); err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	weeklyLogs, err := u.weeklyLogRepository.FindByUserIDAndMonth(input.UserID, input.YearMonth)
	if err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	summarySet := achievementservices.BuildMonthlySummarySetFromWeeklyLogs(input.MonthlyAchievementID, input.UserID, input.YearMonth, weeklyLogs, input.UpdatedAt)
	if err := u.monthlyAchievementRepo.Save(summarySet.Achievement, summarySet.WorkTypeSummaries, summarySet.BreakSummary); err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	return summarySet, nil
}

func (u *summaryUsecase) CreateMonthlySummaryFromWorkLogs(input achievementports.CreateMonthlySummaryInput) (analysismodel.MonthlySummarySet, error) {
	if err := input.Validate(); err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	workLogs, err := u.workLogRepository.FindByUserIDAndMonth(input.UserID, input.YearMonth)
	if err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	breakLogs, err := u.breakLogRepository.FindByUserIDAndMonth(input.UserID, input.YearMonth)
	if err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	summarySet := achievementservices.BuildMonthlySummarySetFromWorkLogs(input.MonthlyAchievementID, input.UserID, input.YearMonth, workLogs, breakLogs, input.UpdatedAt)
	if err := u.monthlyAchievementRepo.Save(summarySet.Achievement, summarySet.WorkTypeSummaries, summarySet.BreakSummary); err != nil {
		return analysismodel.MonthlySummarySet{}, err
	}
	return summarySet, nil
}
