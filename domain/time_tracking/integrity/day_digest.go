// Package integrity provides deterministic representations of time-tracking data.
package integrity

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	analysismodel "time_management/domain/time_tracking/subdomains/analysis/model"
	taskprogressmodel "time_management/domain/time_tracking/subdomains/task_progress/model"
	workrecordmodel "time_management/domain/time_tracking/subdomains/work_record/model"
)

const (
	// CanonicalVersion prefixes every serialized day so future formats can coexist.
	CanonicalVersion = "time-management/day-integrity/v1"
	timestampLayout  = "2006-01-02T15:04:05.000000000Z"
)

var (
	ErrUserIDRequired = errors.New("integrity: user ID is required")
	ErrDayRequired    = errors.New("integrity: day is required")
)

// DayRecords contains the stored records covered by one daily digest.
//
// A day is the half-open UTC interval [00:00:00, 00:00:00 the next day).
// Work logs, daily tasks, and daily achievements belong to the day named by
// WorkDate; break logs belong to the day in which StartedAt occurs. Weekly
// tasks are those whose WeekStartDate is the UTC Monday for the containing
// week. Callers must supply only records selected by those rules.
//
// DailySummarySet is used instead of only DailyAchievement so the persisted
// work-type and break summaries are also protected.
type DayRecords struct {
	UserID            string
	Day               time.Time
	WorkLogs          []workrecordmodel.WorkLog
	BreakLogs         []workrecordmodel.BreakLog
	DailyTasks        []taskprogressmodel.DailyTask
	WeeklyTasks       []taskprogressmodel.WeeklyTask
	DailyAchievements []analysismodel.DailySummarySet
}

// DigestDay returns the SHA-256 digest of SerializeDay(records).
func DigestDay(records DayRecords) ([sha256.Size]byte, error) {
	canonical, err := SerializeDay(records)
	if err != nil {
		return [sha256.Size]byte{}, err
	}
	return sha256.Sum256(canonical), nil
}

// SerializeDay returns the versioned canonical binary representation of records.
//
// Strings and byte sequences are prefixed by an unsigned 64-bit big-endian
// length. Signed integers and durations are 64-bit big-endian values; floats
// use their IEEE-754 bits; optional values have a one-byte presence marker.
// Timestamps are UTC with exactly nanosecond precision. Sections are ordered as
// work logs, break logs, daily tasks, weekly tasks, and daily achievements.
// Within each section records are sorted by ID, with every remaining field in
// the field order below acting as deterministic tie-breakers. Input slices are
// never modified. An empty day serializes the header and five zero counts.
func SerializeDay(records DayRecords) ([]byte, error) {
	if err := validateDay(records); err != nil {
		return nil, err
	}

	var out encoder
	out.rawString(CanonicalVersion)
	out.byte(0)
	out.string(records.Day.UTC().Format("2006-01-02"))
	out.string(records.UserID)
	out.section("work_logs", encodeRecords(records.WorkLogs, encodeWorkLog))
	out.section("break_logs", encodeRecords(records.BreakLogs, encodeBreakLog))
	out.section("daily_tasks", encodeRecords(records.DailyTasks, encodeDailyTask))
	out.section("weekly_tasks", encodeRecords(records.WeeklyTasks, encodeWeeklyTask))
	out.section("daily_achievements", encodeRecords(records.DailyAchievements, encodeDailyAchievement))
	return out.bytes(), nil
}

func validateDay(records DayRecords) error {
	if records.UserID == "" {
		return ErrUserIDRequired
	}
	if records.Day.IsZero() {
		return ErrDayRequired
	}
	day := utcDate(records.Day)
	week := utcWeekStart(day)

	for i, record := range records.WorkLogs {
		if record.UserID != records.UserID {
			return fmt.Errorf("integrity: work log %d has a different user", i)
		}
		if !utcDate(record.WorkDate).Equal(day) {
			return fmt.Errorf("integrity: work log %d is outside the day", i)
		}
	}
	for i, record := range records.BreakLogs {
		if record.UserID != records.UserID {
			return fmt.Errorf("integrity: break log %d has a different user", i)
		}
		if !utcDate(record.StartedAt).Equal(day) {
			return fmt.Errorf("integrity: break log %d is outside the day", i)
		}
	}
	for i, record := range records.DailyTasks {
		if record.UserID != records.UserID {
			return fmt.Errorf("integrity: daily task %d has a different user", i)
		}
		if !utcDate(record.WorkDate).Equal(day) {
			return fmt.Errorf("integrity: daily task %d is outside the day", i)
		}
		if math.IsNaN(record.ProgressRate) || math.IsInf(record.ProgressRate, 0) {
			return fmt.Errorf("integrity: daily task %d has a non-finite progress rate", i)
		}
	}
	for i, record := range records.WeeklyTasks {
		if record.UserID != records.UserID {
			return fmt.Errorf("integrity: weekly task %d has a different user", i)
		}
		if !utcDate(record.WeekStartDate).Equal(week) {
			return fmt.Errorf("integrity: weekly task %d is outside the containing week", i)
		}
	}
	for i, record := range records.DailyAchievements {
		if record.Achievement.UserID != records.UserID {
			return fmt.Errorf("integrity: daily achievement %d has a different user", i)
		}
		if !utcDate(record.Achievement.WorkDate).Equal(day) {
			return fmt.Errorf("integrity: daily achievement %d is outside the day", i)
		}
	}
	return nil
}

func utcDate(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func utcWeekStart(day time.Time) time.Time {
	offset := (int(day.Weekday()) + 6) % 7
	return day.AddDate(0, 0, -offset)
}

func encodeRecords[T any](records []T, encode func(T) []byte) [][]byte {
	encoded := make([][]byte, len(records))
	for i, record := range records {
		encoded[i] = encode(record)
	}
	sort.Slice(encoded, func(i, j int) bool {
		return bytes.Compare(encoded[i], encoded[j]) < 0
	})
	return encoded
}

func encodeWorkLog(record workrecordmodel.WorkLog) []byte {
	var out encoder
	out.string(record.WorkLogID)
	out.string(record.UserID)
	out.string(record.WorkType.WorkTypeID)
	out.string(record.WorkType.WorkTypeName)
	out.time(record.WorkDate)
	out.time(record.StartedAt)
	out.time(record.EndedAt)
	out.duration(record.WorkDuration)
	return out.bytes()
}

func encodeBreakLog(record workrecordmodel.BreakLog) []byte {
	var out encoder
	out.string(record.BreakLogID)
	out.string(record.WorkLogID)
	out.string(record.UserID)
	out.time(record.StartedAt)
	out.time(record.EndedAt)
	out.duration(record.BreakDuration)
	out.optionalTime(record.ResumedAt)
	out.optionalDuration(record.ResumeDuration)
	return out.bytes()
}

func encodeDailyTask(record taskprogressmodel.DailyTask) []byte {
	var out encoder
	out.string(record.DailyTaskID)
	out.string(record.WeeklyTaskID)
	out.string(record.UserID)
	out.time(record.WorkDate)
	out.string(record.ProgressDescription)
	out.float64(record.ProgressRate)
	return out.bytes()
}

func encodeWeeklyTask(record taskprogressmodel.WeeklyTask) []byte {
	var out encoder
	out.string(record.WeeklyTaskID)
	out.string(record.UserID)
	out.string(record.WorkTypeID)
	out.time(record.WeekStartDate)
	out.duration(record.TargetDuration)
	out.string(record.GoalDescription)
	return out.bytes()
}

func encodeDailyAchievement(record analysismodel.DailySummarySet) []byte {
	var out encoder
	achievement := record.Achievement
	out.string(achievement.DailyAchievementID)
	out.string(achievement.UserID)
	out.time(achievement.WorkDate)
	out.duration(achievement.TotalWorkDuration)
	out.duration(achievement.TotalBreakDuration)
	out.int64(int64(achievement.WorkCount))
	out.int64(int64(achievement.BreakCount))
	out.time(achievement.UpdatedAt)

	summaries := encodeRecords(record.WorkTypeSummaries, encodeDailyWorkTypeSummary)
	out.section("work_type_summaries", summaries)
	summary := record.BreakSummary
	out.string(summary.DailyAchievementID)
	out.int64(int64(summary.BreakCount))
	out.duration(summary.BreakDuration)
	out.duration(summary.AverageBreakDuration)
	out.duration(summary.AverageResumeDuration)
	out.duration(summary.LongestResumeDuration)
	return out.bytes()
}

func encodeDailyWorkTypeSummary(summary analysismodel.DailyWorkTypeSummary) []byte {
	var out encoder
	out.string(summary.DailyAchievementID)
	out.string(summary.WorkTypeID)
	out.int64(int64(summary.WorkCount))
	out.duration(summary.WorkDuration)
	out.duration(summary.AverageWorkDuration)
	return out.bytes()
}

type encoder struct {
	bytes.Buffer
}

func (e *encoder) bytes() []byte {
	return append([]byte(nil), e.Bytes()...)
}

func (e *encoder) byte(value byte) {
	_ = e.WriteByte(value)
}

func (e *encoder) rawString(value string) {
	_, _ = e.WriteString(value)
}

func (e *encoder) uint64(value uint64) {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], value)
	_, _ = e.Write(data[:])
}

func (e *encoder) int64(value int64) {
	e.uint64(uint64(value))
}

func (e *encoder) string(value string) {
	e.uint64(uint64(len(value)))
	e.rawString(value)
}

func (e *encoder) data(value []byte) {
	e.uint64(uint64(len(value)))
	_, _ = e.Write(value)
}

func (e *encoder) time(value time.Time) {
	e.string(value.UTC().Format(timestampLayout))
}

func (e *encoder) duration(value time.Duration) {
	e.int64(int64(value))
}

func (e *encoder) float64(value float64) {
	if value == 0 {
		value = 0 // Canonicalize negative zero.
	}
	e.uint64(math.Float64bits(value))
}

func (e *encoder) optionalTime(value *time.Time) {
	if value == nil {
		e.byte(0)
		return
	}
	e.byte(1)
	e.time(*value)
}

func (e *encoder) optionalDuration(value *time.Duration) {
	if value == nil {
		e.byte(0)
		return
	}
	e.byte(1)
	e.duration(*value)
}

func (e *encoder) section(name string, records [][]byte) {
	e.string(name)
	e.uint64(uint64(len(records)))
	for _, record := range records {
		e.data(record)
	}
}
