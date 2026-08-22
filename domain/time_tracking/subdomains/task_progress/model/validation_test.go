package model

import (
	"testing"
	"time"
)

func TestWeeklyTaskValidateRequiresWorkType(t *testing.T) {
	weeklyTask := NewWeeklyTask("00000000-0000-4000-8000-000000000701", "00000000-0000-4000-8000-000000000101", "", day(), time.Hour, "study")

	if err := weeklyTask.Validate(); err != ErrWorkTypeRequired {
		t.Fatalf("expected work type error, got %v", err)
	}
}

func TestDailyTaskValidateRejectsOutOfRangeProgressRate(t *testing.T) {
	dailyTask := NewDailyTask("00000000-0000-4000-8000-000000000801", "00000000-0000-4000-8000-000000000701", "00000000-0000-4000-8000-000000000101", day(), "done", 1.5)

	if err := dailyTask.Validate(); err != ErrProgressRateOutOfRange {
		t.Fatalf("expected progress rate error, got %v", err)
	}
}

func day() time.Time {
	return time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
}
