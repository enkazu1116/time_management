package attestation

import (
	"context"
	"crypto/sha256"
	"errors"
	"sync"
	"time"

	"time_management/domain/time_tracking/integrity"
	achievementrepositories "time_management/domain/time_tracking/subdomains/achievement/interfaces/repositories"
	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressrepositories "time_management/domain/time_tracking/subdomains/task_progress/interfaces/repositories"
	workrecordrepositories "time_management/domain/time_tracking/subdomains/work_record/interfaces/repositories"
)

const (
	StatusMatch       = "match"
	StatusMismatch    = "mismatch"
	StatusNotAttested = "not-attested"
)

var (
	ErrLocationRequired  = errors.New("attestation: business location is required")
	ErrSubjectCommitment = errors.New("attestation: on-chain subject commitment mismatch")
)

type UIDStore interface {
	Save(userID, date, location string, uid [sha256.Size]byte) error
	Latest(userID, date, location string) ([sha256.Size]byte, bool, error)
}

type Service struct {
	writer            *Writer
	workLogs          workrecordrepositories.WorkLogRepository
	breakLogs         workrecordrepositories.BreakLogRepository
	dailyTasks        taskprogressrepositories.DailyTaskRepository
	weeklyTasks       taskprogressrepositories.WeeklyTaskRepository
	dailyAchievements achievementrepositories.DailyAchievementRepository
	uids              UIDStore
	location          *time.Location
	attestMu          sync.Mutex
}

type Result struct {
	Status string
	UID    [sha256.Size]byte
}

func NewService(
	writer *Writer,
	workLogs workrecordrepositories.WorkLogRepository,
	breakLogs workrecordrepositories.BreakLogRepository,
	dailyTasks taskprogressrepositories.DailyTaskRepository,
	weeklyTasks taskprogressrepositories.WeeklyTaskRepository,
	dailyAchievements achievementrepositories.DailyAchievementRepository,
	uids UIDStore,
	location *time.Location,
) (*Service, error) {
	if location == nil {
		return nil, ErrLocationRequired
	}
	return &Service{
		writer: writer, workLogs: workLogs, breakLogs: breakLogs,
		dailyTasks: dailyTasks, weeklyTasks: weeklyTasks,
		dailyAchievements: dailyAchievements, uids: uids, location: location,
	}, nil
}

func (service *Service) Attest(ctx context.Context, userID string, date time.Time) (Result, error) {
	// The demo uses an in-memory locator store, so serialize the complete
	// read-attest-save operation to keep each refUID chain linear.
	service.attestMu.Lock()
	defer service.attestMu.Unlock()

	records, err := service.dayRecords(userID, date)
	if err != nil {
		return Result{}, err
	}
	dateKey := localDate(date, service.location).Format("2006-01-02")
	refUID, _, err := service.uids.Latest(userID, dateKey, service.location.String())
	if err != nil {
		return Result{}, err
	}
	uid, err := service.writer.Attest(ctx, records, refUID)
	if err != nil {
		return Result{}, err
	}
	if err := service.uids.Save(userID, dateKey, service.location.String(), uid); err != nil {
		return Result{}, err
	}
	return Result{Status: StatusMatch, UID: uid}, nil
}

func (service *Service) Verify(ctx context.Context, userID string, date time.Time) (Result, error) {
	dateKey := localDate(date, service.location).Format("2006-01-02")
	uid, found, err := service.uids.Latest(userID, dateKey, service.location.String())
	if err != nil {
		return Result{}, err
	}
	if !found {
		return Result{Status: StatusNotAttested}, nil
	}
	records, err := service.dayRecords(userID, date)
	if err != nil {
		return Result{}, err
	}
	currentDigest, err := integrity.DigestDay(records)
	if err != nil {
		return Result{}, err
	}
	onChain, err := service.writer.Read(ctx, uid)
	if err != nil {
		return Result{}, err
	}
	if onChain.SubjectCommitment != service.writer.SubjectCommitment(records) {
		return Result{}, ErrSubjectCommitment
	}
	status := StatusMismatch
	if onChain.DayDigest == currentDigest {
		status = StatusMatch
	}
	return Result{Status: status, UID: uid}, nil
}

func (service *Service) dayRecords(userID string, date time.Time) (integrity.DayRecords, error) {
	day := localDate(date, service.location)
	dateOnly := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	weekStart := dateOnly.AddDate(0, 0, -((int(dateOnly.Weekday()) + 6) % 7))

	workLogs, err := service.workLogs.FindByUserIDAndDate(userID, dateOnly)
	if err != nil {
		return integrity.DayRecords{}, err
	}
	breakLogs, err := service.breakLogs.FindByUserIDAndDate(userID, day)
	if err != nil {
		return integrity.DayRecords{}, err
	}
	dailyTasks, err := service.dailyTasks.FindByUserIDAndDate(userID, dateOnly)
	if err != nil {
		return integrity.DayRecords{}, err
	}
	weeklyTasks, err := service.weeklyTasks.FindByUserIDAndWeek(userID, weekStart)
	if err != nil {
		return integrity.DayRecords{}, err
	}
	achievement, err := service.dailyAchievements.FindByUserIDAndDate(userID, dateOnly)
	achievements := []analysismodel.DailySummarySet(nil)
	if err == nil {
		achievements = append(achievements, achievement)
	} else if !errors.Is(err, achievementrepositories.ErrNotFound) {
		return integrity.DayRecords{}, err
	}
	return integrity.DayRecords{
		UserID: userID, Day: day, Location: service.location,
		WorkLogs: workLogs, BreakLogs: breakLogs,
		DailyTasks: dailyTasks, WeeklyTasks: weeklyTasks,
		DailyAchievements: achievements,
	}, nil
}

func localDate(date time.Time, location *time.Location) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}
