package memory

import (
	"errors"
	"sync"
	"time"

	achievementmodel "time_management/domain/time_tracking/subdomains/achievement/model"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
	usermodel "time_management/domain/user/model"
	"time_management/infrastructure/util/messages"
)

var ErrNotFound = errors.New(messages.RecordNotFound)

type Store struct {
	mu sync.RWMutex

	users map[string]usermodel.User

	workLogs    map[string]workrecordmodel.WorkLog
	breakLogs   map[string]workrecordmodel.BreakLog
	weeklyTasks map[string]taskprogressmodel.WeeklyTask
	dailyTasks  map[string]taskprogressmodel.DailyTask
	dailyLogs   map[string]taskprogressmodel.DailyLog
	weeklyLogs  map[string]taskprogressmodel.WeeklyLog

	dailySummaries   map[string]analysismodel.DailySummarySet
	monthlySummaries map[string]analysismodel.MonthlySummarySet
}

func NewStore() *Store {
	return &Store{
		users:            map[string]usermodel.User{},
		workLogs:         map[string]workrecordmodel.WorkLog{},
		breakLogs:        map[string]workrecordmodel.BreakLog{},
		weeklyTasks:      map[string]taskprogressmodel.WeeklyTask{},
		dailyTasks:       map[string]taskprogressmodel.DailyTask{},
		dailyLogs:        map[string]taskprogressmodel.DailyLog{},
		weeklyLogs:       map[string]taskprogressmodel.WeeklyLog{},
		dailySummaries:   map[string]analysismodel.DailySummarySet{},
		monthlySummaries: map[string]analysismodel.MonthlySummarySet{},
	}
}

func (s *Store) Create(user usermodel.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.UserID] = user
	return nil
}

func (s *Store) Update(user usermodel.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[user.UserID]; !ok {
		return ErrNotFound
	}
	s.users[user.UserID] = user
	return nil
}

func (s *Store) Delete(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, userID)
	return nil
}

func (s *Store) FindByUserID(userID string) (usermodel.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return usermodel.User{}, ErrNotFound
	}
	return user, nil
}

func (s *Store) CreateWorkLog(workLog workrecordmodel.WorkLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workLogs[workLog.WorkLogID] = workLog
	return nil
}

func (s *Store) UpdateWorkLog(workLog workrecordmodel.WorkLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.workLogs[workLog.WorkLogID]; !ok {
		return ErrNotFound
	}
	s.workLogs[workLog.WorkLogID] = workLog
	return nil
}

func (s *Store) DeleteWorkLog(workLogID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.workLogs, workLogID)
	return nil
}

func (s *Store) FindByWorkLogID(workLogID string) (workrecordmodel.WorkLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	workLog, ok := s.workLogs[workLogID]
	if !ok {
		return workrecordmodel.WorkLog{}, ErrNotFound
	}
	return workLog, nil
}

func (s *Store) FindWorkLogsByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.WorkLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	day := normalizeDate(workDate)
	workLogs := []workrecordmodel.WorkLog{}
	for _, workLog := range s.workLogs {
		if workLog.UserID == userID && sameDate(workLog.WorkDate, day) {
			workLogs = append(workLogs, workLog)
		}
	}
	return workLogs, nil
}

func (s *Store) FindWorkLogsByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.WorkLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeWeekStart(weekStartDate)
	end := start.AddDate(0, 0, 7)
	workLogs := []workrecordmodel.WorkLog{}
	for _, workLog := range s.workLogs {
		if workLog.UserID == userID && inDateRange(workLog.WorkDate, start, end) {
			workLogs = append(workLogs, workLog)
		}
	}
	return workLogs, nil
}

func (s *Store) FindWorkLogsByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.WorkLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeMonth(yearMonth)
	end := start.AddDate(0, 1, 0)
	workLogs := []workrecordmodel.WorkLog{}
	for _, workLog := range s.workLogs {
		if workLog.UserID == userID && inDateRange(workLog.WorkDate, start, end) {
			workLogs = append(workLogs, workLog)
		}
	}
	return workLogs, nil
}

func (s *Store) CreateBreakLog(breakLog workrecordmodel.BreakLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.breakLogs[breakLog.BreakLogID] = breakLog
	return nil
}

func (s *Store) UpdateBreakLog(breakLog workrecordmodel.BreakLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.breakLogs[breakLog.BreakLogID]; !ok {
		return ErrNotFound
	}
	s.breakLogs[breakLog.BreakLogID] = breakLog
	return nil
}

func (s *Store) DeleteBreakLog(breakLogID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.breakLogs, breakLogID)
	return nil
}

func (s *Store) FindByBreakLogID(breakLogID string) (workrecordmodel.BreakLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	breakLog, ok := s.breakLogs[breakLogID]
	if !ok {
		return workrecordmodel.BreakLog{}, ErrNotFound
	}
	return breakLog, nil
}

func (s *Store) FindBreakLogsByUserIDAndDate(userID string, workDate time.Time) ([]workrecordmodel.BreakLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	day := normalizeDate(workDate)
	breakLogs := []workrecordmodel.BreakLog{}
	for _, breakLog := range s.breakLogs {
		if breakLog.UserID == userID && sameDate(breakLog.StartedAt, day) {
			breakLogs = append(breakLogs, breakLog)
		}
	}
	return breakLogs, nil
}

func (s *Store) FindBreakLogsByUserIDAndWeek(userID string, weekStartDate time.Time) ([]workrecordmodel.BreakLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeWeekStart(weekStartDate)
	end := start.AddDate(0, 0, 7)
	breakLogs := []workrecordmodel.BreakLog{}
	for _, breakLog := range s.breakLogs {
		if breakLog.UserID == userID && inDateRange(breakLog.StartedAt, start, end) {
			breakLogs = append(breakLogs, breakLog)
		}
	}
	return breakLogs, nil
}

func (s *Store) FindBreakLogsByUserIDAndMonth(userID string, yearMonth time.Time) ([]workrecordmodel.BreakLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeMonth(yearMonth)
	end := start.AddDate(0, 1, 0)
	breakLogs := []workrecordmodel.BreakLog{}
	for _, breakLog := range s.breakLogs {
		if breakLog.UserID == userID && inDateRange(breakLog.StartedAt, start, end) {
			breakLogs = append(breakLogs, breakLog)
		}
	}
	return breakLogs, nil
}

func (s *Store) CreateWeeklyTask(weeklyTask taskprogressmodel.WeeklyTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weeklyTasks[weeklyTask.WeeklyTaskID] = weeklyTask
	return nil
}

func (s *Store) UpdateWeeklyTask(weeklyTask taskprogressmodel.WeeklyTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.weeklyTasks[weeklyTask.WeeklyTaskID]; !ok {
		return ErrNotFound
	}
	s.weeklyTasks[weeklyTask.WeeklyTaskID] = weeklyTask
	return nil
}

func (s *Store) DeleteWeeklyTask(weeklyTaskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.weeklyTasks, weeklyTaskID)
	return nil
}

func (s *Store) FindByWeeklyTaskID(weeklyTaskID string) (taskprogressmodel.WeeklyTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	weeklyTask, ok := s.weeklyTasks[weeklyTaskID]
	if !ok {
		return taskprogressmodel.WeeklyTask{}, ErrNotFound
	}
	return weeklyTask, nil
}

func (s *Store) FindWeeklyTasksByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.WeeklyTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeWeekStart(weekStartDate)
	weeklyTasks := []taskprogressmodel.WeeklyTask{}
	for _, weeklyTask := range s.weeklyTasks {
		if weeklyTask.UserID == userID && sameDate(weeklyTask.WeekStartDate, start) {
			weeklyTasks = append(weeklyTasks, weeklyTask)
		}
	}
	return weeklyTasks, nil
}

func (s *Store) CreateDailyTask(dailyTask taskprogressmodel.DailyTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dailyTasks[dailyTask.DailyTaskID] = dailyTask
	return nil
}

func (s *Store) UpdateDailyTask(dailyTask taskprogressmodel.DailyTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dailyTasks[dailyTask.DailyTaskID]; !ok {
		return ErrNotFound
	}
	s.dailyTasks[dailyTask.DailyTaskID] = dailyTask
	return nil
}

func (s *Store) DeleteDailyTask(dailyTaskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dailyTasks, dailyTaskID)
	return nil
}

func (s *Store) FindByDailyTaskID(dailyTaskID string) (taskprogressmodel.DailyTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dailyTask, ok := s.dailyTasks[dailyTaskID]
	if !ok {
		return taskprogressmodel.DailyTask{}, ErrNotFound
	}
	return dailyTask, nil
}

func (s *Store) FindDailyTasksByUserIDAndDate(userID string, workDate time.Time) ([]taskprogressmodel.DailyTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	day := normalizeDate(workDate)
	dailyTasks := []taskprogressmodel.DailyTask{}
	for _, dailyTask := range s.dailyTasks {
		if dailyTask.UserID == userID && sameDate(dailyTask.WorkDate, day) {
			dailyTasks = append(dailyTasks, dailyTask)
		}
	}
	return dailyTasks, nil
}

func (s *Store) CreateDailyLog(dailyLog taskprogressmodel.DailyLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dailyLogs[dailyLog.DailyLogID] = dailyLog
	return nil
}

func (s *Store) UpdateDailyLog(dailyLog taskprogressmodel.DailyLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dailyLogs[dailyLog.DailyLogID]; !ok {
		return ErrNotFound
	}
	s.dailyLogs[dailyLog.DailyLogID] = dailyLog
	return nil
}

func (s *Store) DeleteDailyLog(dailyLogID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dailyLogs, dailyLogID)
	return nil
}

func (s *Store) FindByDailyLogID(dailyLogID string) (taskprogressmodel.DailyLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	dailyLog, ok := s.dailyLogs[dailyLogID]
	if !ok {
		return taskprogressmodel.DailyLog{}, ErrNotFound
	}
	return dailyLog, nil
}

func (s *Store) FindDailyLogsByUserIDAndWeek(userID string, weekStartDate time.Time) ([]taskprogressmodel.DailyLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeWeekStart(weekStartDate)
	end := start.AddDate(0, 0, 7)
	dailyLogs := []taskprogressmodel.DailyLog{}
	for _, dailyLog := range s.dailyLogs {
		if dailyLog.UserID == userID && inDateRange(dailyLog.WorkDate, start, end) {
			dailyLogs = append(dailyLogs, dailyLog)
		}
	}
	return dailyLogs, nil
}

func (s *Store) CreateWeeklyLog(weeklyLog taskprogressmodel.WeeklyLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weeklyLogs[weeklyLog.WeeklyLogID] = weeklyLog
	return nil
}

func (s *Store) UpdateWeeklyLog(weeklyLog taskprogressmodel.WeeklyLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.weeklyLogs[weeklyLog.WeeklyLogID]; !ok {
		return ErrNotFound
	}
	s.weeklyLogs[weeklyLog.WeeklyLogID] = weeklyLog
	return nil
}

func (s *Store) DeleteWeeklyLog(weeklyLogID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.weeklyLogs, weeklyLogID)
	return nil
}

func (s *Store) SaveAll(weeklyLogs []taskprogressmodel.WeeklyLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, weeklyLog := range weeklyLogs {
		s.weeklyLogs[weeklyLog.WeeklyLogID] = weeklyLog
	}
	return nil
}

func (s *Store) FindByWeeklyLogID(weeklyLogID string) (taskprogressmodel.WeeklyLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	weeklyLog, ok := s.weeklyLogs[weeklyLogID]
	if !ok {
		return taskprogressmodel.WeeklyLog{}, ErrNotFound
	}
	return weeklyLog, nil
}

func (s *Store) FindWeeklyLogsByUserIDAndMonth(userID string, yearMonth time.Time) ([]taskprogressmodel.WeeklyLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	start := normalizeMonth(yearMonth)
	end := start.AddDate(0, 1, 0)
	weeklyLogs := []taskprogressmodel.WeeklyLog{}
	for _, weeklyLog := range s.weeklyLogs {
		if weeklyLog.UserID == userID && inDateRange(weeklyLog.WeekStartDate, start, end) {
			weeklyLogs = append(weeklyLogs, weeklyLog)
		}
	}
	return weeklyLogs, nil
}

func (s *Store) SaveDaily(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dailySummaries[achievement.DailyAchievementID] = analysismodel.DailySummarySet{
		Achievement:       achievement,
		WorkTypeSummaries: workTypeSummaries,
		BreakSummary:      breakSummary,
	}
	return nil
}

func (s *Store) UpdateDaily(achievement achievementmodel.DailyAchievement, workTypeSummaries []analysismodel.DailyWorkTypeSummary, breakSummary analysismodel.DailyBreakSummary) error {
	return s.SaveDaily(achievement, workTypeSummaries, breakSummary)
}

func (s *Store) DeleteDaily(dailyAchievementID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dailySummaries, dailyAchievementID)
	return nil
}

func (s *Store) FindDailyByUserIDAndDate(userID string, workDate time.Time) (analysismodel.DailySummarySet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	day := normalizeDate(workDate)
	for _, summarySet := range s.dailySummaries {
		if summarySet.Achievement.UserID == userID && sameDate(summarySet.Achievement.WorkDate, day) {
			return summarySet, nil
		}
	}
	return analysismodel.DailySummarySet{}, ErrNotFound
}

func (s *Store) SaveMonthly(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.monthlySummaries[achievement.MonthlyAchievementID] = analysismodel.MonthlySummarySet{
		Achievement:       achievement,
		WorkTypeSummaries: workTypeSummaries,
		BreakSummary:      breakSummary,
	}
	return nil
}

func (s *Store) UpdateMonthly(achievement achievementmodel.MonthlyAchievement, workTypeSummaries []analysismodel.MonthlyWorkTypeSummary, breakSummary analysismodel.MonthlyBreakSummary) error {
	return s.SaveMonthly(achievement, workTypeSummaries, breakSummary)
}

func (s *Store) DeleteMonthly(monthlyAchievementID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.monthlySummaries, monthlyAchievementID)
	return nil
}

func (s *Store) FindMonthlyByUserIDAndMonth(userID string, yearMonth time.Time) (analysismodel.MonthlySummarySet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	month := normalizeMonth(yearMonth)
	for _, summarySet := range s.monthlySummaries {
		if summarySet.Achievement.UserID == userID && sameDate(summarySet.Achievement.YearMonth, month) {
			return summarySet, nil
		}
	}
	return analysismodel.MonthlySummarySet{}, ErrNotFound
}

func normalizeDate(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func normalizeMonth(value time.Time) time.Time {
	year, month, _ := value.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, value.Location())
}

func normalizeWeekStart(value time.Time) time.Time {
	day := normalizeDate(value)
	offset := (int(day.Weekday()) + 6) % 7
	return day.AddDate(0, 0, -offset)
}

func sameDate(left time.Time, right time.Time) bool {
	return normalizeDate(left).Equal(normalizeDate(right))
}

func inDateRange(value time.Time, start time.Time, end time.Time) bool {
	day := normalizeDate(value)
	return !day.Before(start) && day.Before(end)
}
